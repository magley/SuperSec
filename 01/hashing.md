# Heširanje

## Zadatak

Dizajnirati mehanizam hešovanja sa ciljem da se zaštiti poverljivost (confidentiality) korisničkih lozinki.

1) Istražiti različite algoritme i odabrati najbezbedniji;
2) Ispitati konfiguracione parametre odabranog algoritma, i otkriti koja bi to bila preporučena praksa za
konfiguraciju;
3) Odabrati pouzdanog provajdera;
4) Istražiti da li poslednja verzija za implementaciju ima ozbiljnijih ranjivosti;
5) Specificirati zahteve za bezbednu implementaciju heš mehanizma koristeći sve do sada nabrojano.

### 1. Razmatrani algoritmi

Algoritmi za heširanje se mogu grubo podeliti u dve grupe: brze i spore.

Brzi algoritmi se primarno koriste za autentifikaciju poruka i nisu pogodni za heširanje poverljivih informacija poput lozinki. Razlog tome je što napadač pri "krekovanju" lozinke hešira ogroman broj lozinki i poredi sa heš vrednosti tražene lozinke, te je jedna mera zaštite od ovakvog napada produživanje vremena potrebnog za heširanje. Primeri ovakvih algoritama su MD5 i SHA-1, i njih nećemo dalje razmatrati.

Interesantniji su nam spori algoritmi koji su upravo namenjeni heširanju poverljivih informacija. Po redosledu preporuke nabrajamo:
- Argon2id - pobednik na [takmičenju heširanja lozinki 2015. godine](https://en.wikipedia.org/wiki/Password_Hashing_Competition)
- scrypt - preporučen kada Argon2id nije dostupan
- bcrypt - preporučen u legacy sistemima
- PBKDF2 - preporučen od strane NIST-a i ima FIPS-140 validirane implementacije, te se preferira ako su nam ovo zahtevi

U daljem razmatranju fokusiraćemo se na Argon2id algoritam. Skrećemo pažnju da je Argon2id samo jedna varijanta Argon2 algoritma, te ćemo u zavisnosti od konteksta opisivati nekada Argon2, a nekada Argon2id.

### 2. Konfiguracioni parametri

Argon2id ima tri konfiguraciona parametra:
- m - bazni minimum minimalne veličine memorije (base minimum of the minimum memory size)
- t - minimalni broj iteracija
- p - stepen paralelizma

Slede preporučene (i bezbednosno ekvivalentne) vrednosti parametara:
- m=47104 (46 MiB), t=1, p=1 (Ne koristiti sa Argon2i)
- m=19456 (19 MiB), t=2, p=1 (Ne koristiti sa Argon2i)
- m=12288 (12 MiB), t=3, p=1
- m=9216 (9 MiB), t=4, p=1
- m=7168 (7 MiB), t=5, p=1

### 3. Provajder algoritma

Za provajdera algoritma odabrali smo Spring-ovu Spring Security Crypto biblioteku. Spring spada u najpopularnije veb frejmvorke i koristi se već duži niz godina, što je dobar znak njegove pouzdanosti.

U projektu koji koristi Maven na sledeći način dodajemo Argon2 implementaciju:
```xml
<dependency>
    <groupId>org.springframework.security</groupId>
    <artifactId>spring-security-crypto</artifactId>
    <version>6.0.3</version>
</dependency>
```

Primer upotrebe:
```java
    String rawPassword = "password";
    Argon2PasswordEncoder arg2SpringSecurity = new Argon2PasswordEncoder(saltLength, hashLength, parallelism, memory, iterations);
    String springBouncyHash = arg2SpringSecurity.encode(rawPassword);
```

### 4. Ranjivosti

Ne izgleda da je Argon2id algoritam izložen opasnim ranjivostima. Naišli smo na diskusiju o ranjivostima Argon2i algoritma i njihovom mogućem prenošenju na Argon2id, ali se donosi zaključak da ovo nije moguće.

U dokumentaciji `Argon2PasswordEncoder` klase spominje se da implementacija ne vrši eksploataciju paralelizma/optimizacija koje možda radi napadač, te da postoji nepotrebna asimetrija između napadača i branioca.

Nije pronađena nijedna aktuelna ranjivost u SNYK i CVE bazama ranjivosti.

### 5. Zahtevi bezbedne implementacije heš mehanizma

Iz svega navedenog možemo zaključiti sledeće zahteve:
- Heš algoritam mora biti spor, trošeći što više procesorkih ciklusa i memorije
- Parametri algoritma (poput upravo spomenute brzine) moraju biti konfigurabilni kako bi krajnji korisnik mogao algoritam konfigurisati svojim potrebama (odnos potrošnje procesorskih ciklusa i memorije, opšta brzina...)
 - Algoritam mora podržavati *salting*

Pored ovih zahteva specifičnih za heširanje lozinki, nabrajamo i opšte zahteve:
- Za dva slična ulaza imamo značajno različite izlaze
- Jednosmernost - nije moguće preko izlaza izračunati ulaz
- Teško je proizvesti koliziju (različiti ulazi ne smeju dati isti izlaz)

## Reference

[OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)

[Hashing with Argon2 in Java](https://www.baeldung.com/java-argon2-hashing)

[Java Spring (Boot) popularity](https://www.statista.com/statistics/1124699/worldwide-developer-survey-most-used-frameworks-web/)

[Argon2id safety discussion](https://crypto.stackexchange.com/questions/80379/is-argon2id-safe-from-the-discovered-argon2i-issues)

[Argon2PasswordEncoder Attacker/defender assymetry note](https://docs.spring.io/spring-security/site/docs/current/api/org/springframework/security/crypto/argon2/Argon2PasswordEncoder.html)

[Basic requirements of hashing functions](https://www.quora.com/What-are-the-6-requirements-of-a-secure-hashing-function)