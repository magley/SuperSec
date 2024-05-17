# 4. Nigerian Prince

Flag: `UNS{EM4IL_5P4M_AG4N?}`

Čitajući `email.txt` iz teksta se kao termini za pretragu ističu dva termina:
- Senate bill 1622 , Title 8 ; Section 308
- Senate bill 1916 , Title 1 ; Section 305

Pretragom ovih termina nailazimo na [forum stranicu](https://www.vbforums.com/showthread.php?157720-Decode-the-SPAM) koja daje instrukcije za dekodiranje spam poruka i upućuje nas na https://www.spammimic.com/decode.shtml. Unošenjem teksta iz `email.txt` fajla i dekodiranjem dobijamo traženi flag.

# 5. Educational Purposes Only

Flag: `UNS{V3RY_OLD_4RCH1V3}`

## 1. Date when Faculty of Technical Sciences officialy opened

Na [wikipedia](https://en.wikipedia.org/wiki/University_of_Novi_Sad_Faculty_of_Technical_Sciences) stranici fakulteta nalazimo da je fakultet osnovan 18. maja 1960. godine.

## 2. First name of the person who held the position of dean at the faculty from 01.10.1975. until September 30, 1977.

Pretragom termina "ftn dekani" nailazimo na stranicu o [istorijatu funkcije dekana](http://www.ftn.uns.ac.rs/n508315396/istorijat-funkcije-dekan). Za navedeni vremenski opseg iz tabele na stranici vidmo da je dekan bio Dragutin Zelenović.

## 3. The date when the FTN website was launched. (Date Format : DD/MM/YYYY)

Pokušano je nekoliko pristupa:

1. U _Page Source_ stranice pokušano je izvući informaciju o tome kada je sajt prvi put objavljen, pretragom termina date/published/datePublished

2. Na [wayback machine](https://web.archive.org/web/20090501000000*/ftn.uns.ac.rs) pronađen je najstariji snapshot sajta

3. Pokušano je nekoliko sajtova za proveru starosti sajta ([1](https://websiteage.org/), [2](https://www.duplichecker.com/domain-age-checker.php), [3](https://smallseotools.com/domain-age-checker/), [4](https://carbondate.cs.odu.edu/))

4. Konačno je pretragom termina "objavljen sajt fakulteta tehnickih nauka" pronađena [sledeća](http://www.ftn.uns.ac.rs/102746447/novi-sajt) stranica na kojoj nalazimo _blogpost_ o novom sajtu datiran na 18.05.2005.

## 4. The year when studies in the field of "Poštanski saobraćaj i telekomunikacije" were introduced.

Na [stranici](http://www.ftn.uns.ac.rs/1460455804/postanski-saobracaj-i-telekomunikacije) o kursu istraženo je nekoliko linkova: [akreditacija](http://www.ftn.uns.ac.rs/1063714277/postanski-saobracaj-i-telekomunikacije), [opis kursa](http://www.ftn.uns.ac.rs/1063714277/postanski-saobracaj-i-telekomunikacije), kao i nastavni planovi kursa. Na stranicama je sa Ctrl+F tražen tekst koji počinje sa "19" ili "20", odnosno moguće prve cifre godine početka programa.

Najkorisnije informacije su nađene na stranicama o nastavnim planovima kursa ([2005](http://www.ftn.uns.ac.rs/1167044656/), [2007](http://www.ftn.uns.ac.rs/2054548337/), [2009](http://www.ftn.uns.ac.rs/n1950895790/)...). Na svim ovim stranicama postoji tabela u kojoj piše da je "Godina u kojoj je započeta realizacija studijskog programa" 2005. godina. Međutim, poredeći MD5 heš sa traženim vidimo da ovo nije očekivano rešenje.

_Brute force_ pogađanjem nalazimo da je tražena godina 1999. U pokušaju pronalaska izvora ove informacije pretražujemo "postanski saobracaj ftn "1999"", međutim ne nalazimo izvor.

## Konačna šifra

Konačna šifra za otvaranje rar arhive je `18/05/1960Dragutin18/05/20051999`. Unutar arhive je `flag.png` fajl sa kojeg vidimo da je traženi flag `UNS{V3RY_OLD_4RCH1V3}`