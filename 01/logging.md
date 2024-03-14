# Logovanje

## Zadatak

Dizajnirati mehanizam za logovanje događaja koji odgovara na sledeće zahteve:

1) Log datoteke moraju pružiti informacije potrebne za razrešavanje problema;
2) Svi događaji za koje su akteri bitni moraju biti zapisani, sa dovoljno informacija kako akteri ne bi mogli da
poriču odgovornost (non-repudiation). Potrebno je obezbediti lako izdvajanje tih događaja;
3) Stavke log datoteke ne smeju sadržati osetljive podatke;
4) Mehanizam za logovanje mora biti pouzdan, mora obezbediti dostupnost i integritet log datoteka;
5) Stavke log datoteke moraju precizno iskazati vreme nastanka;
6) Mehanizam za logovanje mora stremiti ka tome da su logovi uredni, da je “pretrpanost” minimalizovana.

Zadatak je istražiti kako se svaki od ovih zahteva može ispuniti, i specificirati konkretne korake implementacije za dizajnirani mehanizam. Korišćenje konkretnih rešenja je dozvoljeno, ako dato rešenje ispunjava sve zahteve, ili se može proširiti tako da ispunjava.

## Rešenje

Biće prikazano u Javi.

### Konkretno rešenje

Java nudi ugrađenu podršku za logove putem `java.util.logging` paketa (**JUL**). Za male projekte gde logovanje nije važno, JUL je dovoljan. U suprotnom (ako je potrebno: parametrizacija, filtriranje, rotacija; i da ove stvari rade brzo), preporuka je da se koristi moćniji API - **SLF4J**. 

SLF4J je samo fasada za funkcionalnosti logovanja. Da bismo ga koristili, treba nam biblioteka koja implementira SLF4J. U praksi se najviše koriste **Log4j** i **Logback**. 

Odabrali smo **Logback** jer nudi dodatne funkcionalnosti za arhiviranje logova, brži je i Spring podrazumevano koristi logback.

Dodavanje Logback-a:

```xml
<!-- Ako se ne koristi Spring (logback-classic će uvući slf4j) -->
<dependency> 
    <groupId>ch.qos.logback</groupId> 
    <artifactId>logback-classic</artifactId> 
    <version>1.4.5</version> 
</dependency>

<!-- Ako se koristi Spring (starter-web će uvući i logback) -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
</dependency>
```

Konfigurisanje logback-a se radi kroz `logback.xml` odnosno `logback-spring.xml`.

### 1. Logovanje debug informacija

> Log datoteke moraju pružiti informacije potrebne za razrešavanje problema

SLF4J nudi nekoliko nivoa logova. Za debagovanje se koriste `DEBUG` i `TRACE`, dva nivoa najnižeg prioriteta.

`TRACE` biramo kada nam treba finije logovanje u sklopu jednog "procesa" ili funkcije. Za logovanje korak-po-korak, `TRACE` je bolja opcija. U ostalim slučajevima, koristimo `DEBUG`. 

Kada se desi izuzetak, logger bi trebalo da ispiše _stack trace_. Logback ovo podrazumevano radi.

### 2. Neporecivost

> Svi događaji za koje su akteri bitni moraju biti zapisani, sa dovoljno informacija kako akteri ne bi mogli da
poriču odgovornost (non-repudiation). Potrebno je obezbediti lako izdvajanje tih događaja;

Ako koristimo Spring Boot, moguće je dobaviti korisnika koji je povezan sa niti u kojoj se izvršava zahtev koristeći:
```java
TheUser currentUser = (TheUser) SecurityContextHolder.getContext().authentication.principal
```

Zatim na svaki log možemo dodati podatak o korisniku:

```java
MoneyTransaction t = ...; 
logger.info("User " + currentUser.getEmail() + " has deposited " + t.amount + " into his account.");
```

Međutim, ovakav pristup je sklon greškama i propustima. Ovo bismo mogli rešiti tako što koristimo *key-value* mapu u koju na početku obrade zahteva dodamo podatak o npr. korisniku koji šalje zahtev. Logger metode bismo proširili koristeći wrapper koji na zadati string dodaje ključ-vrednost parove iz te mape.

Ovaj mehanizam već postoji u SFL4J-u, i zove se **mapped diagnostic context (MDC)**. Ideja je da logger čuva dodatni kontekst u vidu mape koji se umeće u delove loga.

```java
@PUT
ResponseEntity<*> depositMoney(MoneyTransaction t) {
    TheUser currentUser = ...;

    MDC.put("user.email", currentUser.getEmail());
    logger.info("Deposited " + t.amount + " into own account.");
    MDC.clear();
}
```

```xml
<pattern>
%-4r [%t] %5p %c{1} - %m - user=%X{user.email} %n
<!--                       ^^^^^^^^^^^^^^^^^^^ -->
</pattern>
```

Na ovaj način smo obezbedili neporecivost događaja. Ako želimo kasnije da pretražimo događaj po korisniku, mogli bismo izvršiti filtriranje npr. `user=bob123@gmail.com`.

MDC ima i svoje probleme. Ako se desi greška, logback neće to prikazati kao izuzetak. MDC se vezuje na niti, ali nit ne preuzima kontekst svog roditelja. Ako se izvršavanje sprovodi između dveju ili više niti, MDC neće prebaciti kontekst.

Alternativa MDC-u je strukturirano logovanje, koje je opisano u sekciji "6. Uredni logovi".

### 3. Osetljivi podaci

> Stavke log datoteke ne smeju sadržati osetljive podatke;

Prva linija odbrane bi bila da se **ne loguju osetljivi podaci**. Drugim rečima, treba ručno da pazimo šta logujemo. Šta podrazumeva osetljive podatke zavisi od slučaja korišćenja, prirode podataka, politike korišćenja i zakona države. 

Loggeri automatski pozivaju `.toString()` metodu nad svim objektima prilikom popunjavanja parametrizovanih logova. Ako naivno prosleđujemo objekat, potencijalno može doći do curenja njegovih osetljivih podataka u log.

Gde god je moguće, trebalo bi izbeći plain-text tajne. Šifre se mogu **heširati**, i upotrebom bezbednog algoritma (npr. Argon2 sa dobrim parametrima) nam se skoro garantuje da je nemoguće zloupotrebiti taj podatak ako se pojavi u logovima. Nije moguće heširati sve tajne, npr. osetljivi podaci za kreditnu karticu. U tom slučaju možemo raditi **enkripciju**, ali ovo nas ne štiti u potpunosti.

Umesto logovanja tajne, možemo logovati neki token vezan za tu tajnu, npr. ID korisnika u bazi ili surogatni ključ, umesto ličnih podataka korisnika.

Osetljive podatke ne treba slati u sklopu web API-ja, jer server loguje svaki zahtev i ti podaci se čuvaju kao deo URL-a.

Kredencijali i platni detalji ne smeju da se pojave u logovima. Ovi podaci u nekoj meri prate šablon (ili bismo mogli da ih prisilmo da prate šablon), te bi logger mogao automatski da prepozna takav šablon i maskira odgovarajuće podatke. SLF4J nudi podršku za automatsko sakrivanje podataka na osnovu šablona. Na primer, ako u logu ispišemo tekst oblika: "creditCard=123456" ili ako ispisujemo objekat `{'creditCard': "123456", ...}`, onda bi SLF4j automatski mogao da taj tekst zameni sa zvezdicama (`********`):

```xml
<!-- https://stackoverflow.com/a/61360391 -->
<appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
<encoder class="ch.qos.logback.core.encoder.LayoutWrappingEncoder">
    <layout class="com.bpcbt.micro.utils.PatternMaskingLayout">
        <maskPattern>creditCard=\d+</maskPattern>
        <!-- ^^^^ Ovde se maskiraju podaci na osnovu šablona-->
        <pattern>%d{dd/MM/yyyy HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n%ex</pattern>-->
    </layout>
</encoder>
```

### 4. Pouzdanost, dostupnost i integritet

> Mehanizam za logovanje mora biti pouzdan, mora obezbediti dostupnost i integritet log datoteka;

Logback nije **pouzdan**. Da bi postigao visoke performanse u situacijama kada se loguje velik broj događaja, logback će potencijalno ignorisati logovanje događaja čiji je nivo ispod `WARNING`.

Postoji nekoliko napada koji narušavaju **dostupnost** logova:

1. Nadapač šalje velik broj zahteva što prouzrokuje velikom broju logova. Disk se popunjava logovima što sprečava kako naknadno logovanje, tako i rast drugih vrsta podataka na disku (resursi, baza podataka).

2. Velik broj logova u jedinici vremena može da naruši performanse preostalih servisa koji rade na čvoru gde se obavlja logovanje.

Napad na prostor se rešava centralizacijom: rutiranje svih logova ka jednom čvoru u serveru. Za ovo možemo koristiti syslog daemone poput *rsyslog* ili on-premise rešenja kao *Logstash*. Moguće je izgubiti log datoteke u tranzitu (pogotovo ako se logovi ne čuvaju van centralizovanog servisa barem privremeno). Stoga je poželjno uvesti replikaciju, ali to dovodi do narušavanja konzistentnosti.
Logback je optimizovan za high load scenarije ali to dolazi po cenu pouzdanosti.

**Integritet** obezbeđujemo na nekoliko načina.
Potpisivanjem logova (ili dela logova) bismo učinili naš log fajl _tamper proof_ spolja. Ekstremni slučaj ovoga bi bio _blockchain_-oliki pristup heširanja elemenata na osnovu prethodnih, ali u praksi ovakav pristup je nedovoljno efikasan te se ne primenjuje.

Dodatna zaštita bi podrazumevala enkripciju logova pre pisanja u disk.
Ključ za dekripciju bi bio dostupan samo određenim korisnicima/ulogama u sistemu.

Ako šaljemo logove u centralizovan sistem, neophodno je to raditi preko bezbednog protokola (HTTPs).

Drugi pristup zaštite integriteta bi bio da se logovi repliciraju na više različitih tačaka (neparan broj). Ako postoji nekonzistentnost među replikama, mogli bismo koristiti algoritam za postizanje koncenzusa na osnovu glasanja.

Čuvanje logova na _write-once-read-many_ uređaj garantuje integritet.

### 5. Vreme nastanka

> Stavke log datoteke moraju precizno iskazati vreme nastanka;

Logback nudi podršku za podešavanje formata logova. Jedna od stvari koje se mogu podesiti je vreme nastanka:

```xml
<!-- logback.xml -->

<configuration>
  <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
    <encoder>
      <pattern>%d{yyyy-MM-dd_HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
    </encoder>
  </appender>

  <root level="debug">
    <appender-ref ref="STDOUT" />
  </root>
</configuration>
```

Ako logujemo sa više tačaka u sistemu, potrebno je uskladiti njihove časovnike. Ovo je dužnost samih računara i najbolje se sprovodi kroz lokalni **network time protocol server**.

Preporučuje se čuvanje vremena nastanka svih logova po istom standardu, npr. UTC.

### 6. Uredni logovi

> Mehanizam za logovanje mora stremiti ka tome da su logovi uredni, da je “pretrpanost” minimalizovana.

Na nama je da koristimo minimalan broj reči kada zapisujemo logove. Za deskriptore objekata je najbolje koristiti njihov ID (npr. umesto punog imena korisnika korsititi primarni ključ tog korisnika iz baze podataka).

Upotrebom **strukturiranog logovanja** možemo povećati čitljivost logova. Ispis logova je onda JSON, XML, YAML i sl. umesto linije u plain text-u. Ovo olakšava i računarsku obradu i analizu logova jer logovi imaju strukturu u vidu polja.

Logback podržava strukturirano logovanje putem Fluent API-ja i JSON formatter-a:

```java
// Običan API.
logger.debug("old_T={} new_T={} Temperature change", newT, oldT);

// Fluent API.
logger
    .atDebug()
    .addKeyValue("old_T", oldT)
    .addKeyValue("new_T", newT)
    .log("Temperature change");
```

Ako se opredelimo za logove bez šeme, potrebno je formatirati ih tako da su delovi logova jasno odvojeni i da imaju fiksan razmak.

Sledeći log ispis je teži za čitanje jer su svi delovi odvojeni jednim razmakom.

```
2018-07-29 21:10:29.178 thread-1 INFO com.example.MyService Service started in 3434 ms.
2018-07-29 21:10:29.178 main WARN some.external.Configuration parameter 'foo' is missing. Using default value 'bar'!
2018-07-29 21:10:29.178 scheduler ERROR com.example.jobs.ScheduledJob Scheduled job cancelled due to NullPointerException!
```

Sledeći log ispis je mnogo uredniji.

```
2018-07-29 | 21:10:29.178 | thread-1  | INFO  | com.example.MyService         | Service started in 3434 ms.
2018-07-29 | 21:10:29.178 | main      | WARN  | some.external.Configuration   | Parameter 'foo' is missing. Using default value 'bar'!
2018-07-29 | 21:10:29.178 | scheduler | ERROR | com.example.jobs.ScheduledJob | Scheduled job cancelled due to NullPointerException!
```

Format logova koje Logback ispisuje može se promeniti u konfiguracionom fajlu.

```xml
<!-- logback.xml -->

<configuration>
  <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
    <encoder>
      <pattern>%d{yyyy-MM-dd_HHĐ} | %d{:mm:ss.SSS} | %-10thread | %-5level | %logger{36} | %msg%n</pattern>
    </encoder>
  </appender>
</configuration>
```

Ako u logove ubacujemo nepotrebne podatke, logovi će biti pretrpani.
Minimalan skup podataka koje treba logovati čine:
- vreme i datum logova (prateći neki standard poput ISO-8601)
- nivo ozbiljnosti logova
- nit u kojem se log desio
- loger koji je ispisao log
- log poruka
- stack trace u slučaju greške

Kada odaberemo jedan format logovanja, njega treba da se pridržavamo. 

**Rotiranje logova** podrazumeva periodično arhiviranje i brisanje starih logova radi uštede prostora.
Da bi čitanje logova bilo lakše, treba podesiti odgovarajuć period za rotiranje.
Logback nudi podršku za rotiranje logova.

Pored rotiranja, zgodno je i organizovati logove u više fajlova.
Kreiranje direktorijuma po datumima, podela fajlova po nivou ozbiljnosti su neki od
pristupa koji olakšavaju navigaciju prilikom debagovanja.

## Reference

[Debug logging](https://www.crowdstrike.com/cybersecurity-101/observability/debug-logging/)

[Logging in Java, best practices (Nov 2023)](https://betterstack.com/community/guides/logging/how-to-start-logging-with-java/)

[Reasons to use logback instead of log4j](https://logback.qos.ch/reasonsToSwitch.html)

[Log4j vs Logback vs SLF4J](https://sematext.com/blog/java-logging-frameworks/)

[Spring boot logback](https://www.baeldung.com/spring-boot-logback-log4j2)

[MDC](https://www.baeldung.com/mdc-in-log4j-2-logback)

[OWASP Cheat Sheet for logging](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)

[Logging secrets](https://stackoverflow.com/questions/33671027/logging-security-considerations-and-sensitive-data)

[7 best practices for logging](https://medium.com/@joecrobak/seven-best-practices-for-keeping-sensitive-data-out-of-logs-3d7bbd12904)

[OWASP logging guide](https://owasp.org/www-pdf-archive/OWASP_Logging_Guide.pdf)

[OWASP security logging](https://owasp.org/www-project-proactive-controls/v3/en/c9-security-logging)

[Reliable logging](https://stackoverflow.com/questions/3430894/what-exactly-is-reliable-logging)

[Logback is not reliable](https://stackoverflow.com/questions/7681498/logback-reliability)

[Centralized logging](https://www.loggly.com/ultimate-guide/centralizing-java-logs/)

[Case against MDC (last paragraph)](https://tersesystems.github.io/terse-logback/structured-logging/)

[Readable logs](https://reflectoring.io/logging-format/)

[Log formatting best practices](https://www.sentinelone.com/blog/log-formatting-best-practices-readable/)

[Common Log Format](https://en.wikipedia.org/wiki/Common_Log_Format)

[Log integrity](https://docs.logsentinel.com/advanced/log-integrity/)

[NIST publication for logging](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-92.pdf)

[OWASP ASVS v4.0.3](https://github.com/OWASP/ASVS/tree/v4.0.3#latest-stable-version---403)