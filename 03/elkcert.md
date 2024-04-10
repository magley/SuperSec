## 1)

**Naziv izazova u aplikaciji**: **Login Jim**

**Klasa napada**: Injection

**Uticaj**: Napadač može neovlašćeno pristupiti administratorskom, ili bilo kom drugom nalogu.

**Ranjivosti**: Upotrebom SQL injekcije napadač može zaobići potrebu za unošenjem lozinke naloga.

**Kontramere**: _Prepared Statements_ sa parametrizovanim upitima. Oni nam omogućuju da jasno razlikujemo izvršivi kod od parametra upita.

**Beleške**:

Na kraj polja za unos mejla dadato je `'--`, čime smo postigli efekat da u SQL upitu zakomentarišemo deo upita koji mečuje lozinku.

Možemo zamisliti da upit izgleda ovako:

```sql
SELECT * FROM USERS WHERE email = 'jim@juice-sh.op'--' AND password = ''
```

## 2)

**Naziv izazova u aplikaciji**: **Forged Feedback**

**Klasa napada**: Broken Access Control

**Uticaj**: Napadač može obavljati radnje u ime drugog korisnika, ili obaviti radnje koje njemu samom ne bi trebalo biti autorizovane.

**Ranjivosti**: Modifikacijom HTML forme na klijentskoj strani ili presretanjem HTTP zahteva možemo promeniti ID korisnika.

**Kontramere**:

1. Neizlaganje nepotrebnih informacija o korisniku poput ID-ja
2. Autorizaciju obavljati na bekendu, ne na frontendu

**Beleške**:

Na formi za _Customer Feedback_ pronađeno je nevidljivo ID polje. Nakon što je uklonjen atribut `hidden=""` u polje je moguće uneti ID proizvoljnog korisnika.

## 3)

**Naziv izazova u aplikaciji**: **Admin Registration**

**Klasa napada**: Improper Input Validation

**Uticaj**: Napadač može registrovati korisnika sa ulogom administratora.

**Ranjivosti**: Prilikom registracije korisnika, moguće je proslediti ulogu korisnika kao parametar.

**Kontramere**:

https://cheatsheetseries.owasp.org/cheatsheets/Mass_Assignment_Cheat_Sheet.html

1. Bajndovanje polja na DTO objekte koji sadrže samo tražene atribute
2. _Allow-listing_, _Block-listing_

**Beleške**:

Nakon slanja običnog zahteva za registraciju, u _Network_ tabu _Firefox_-a sam pokušao _Edit and Resend_ zahteva sa dodatim poljem `"type": "admin"`.
Ovo nije postiglo željeni rezultat. Sledeće sam iz JWT tokena ulogovanog administratora zaključio da je traženi oblik bio `"role": "admin"`.
Isto sam mogao zaključiti iz HTTP odgovora zahteva za registraciju. Sledeći pokušaj registracije sa dodatim poljem je bio uspešan.

## 4)

**Naziv izazova u aplikaciji**: **API-only XSS**

**Klasa napada**: XSS

**Uticaj**: Moguće je na bekendu skladištiti malicioznu skriptu koja će se izvršiti na svakom računaru kojem se prikaže.

**Ranjivosti**: Nije izvršena sanitacija HTML-a. Nije validiran unos.

**Kontramere**:

https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html

1. HTML sanitacija
2. Validacija unosa

**Beleške**:

_Hint_ stranica nas naslućuje da se fokusiramo na pregled proizvoda. Eskperimentisanjem možemo zaključiti da je dostupan endpoint `/api/Products/${id}`
koji nam dostavlja podatke o pojedinačnom proizvodu. Na ovaj endpoint je moguće uraditi PUT i u sadržaj HTTP zahteva staviti sledeće:

```json
{"description":"<iframe src=\"javascript:alert(`xss`)\">"}
```

Takođe je bitno u zahtevu postaviti  `Content-Type: application/json`

## 5)

**Naziv izazova u aplikaciji**: **Login Amy**

**Klasa napada** Sensitive Data Exposure

**Uticaj**: Pruža polaznu tačku za druge napade, daje informacije napadaču kojima se može predstaviti kao neko drugi

**Ranjivosti**: Nedostatak dvofaktorske autentifikacije, nedostatak email potvrde korisnika pri pristupu sa nepoznate IP adrese

**Kontramere**:

1. Dvofaktorska autentifikacija
2. Zahtevati email potvrdu korisnika pri pristupu sa nepoznate IP adrese
3. Za klijenta: koristiti nasumične lozinke koje se ne mogu povezati sa klijentom OSINT-om (password manager za skladištenje i generisanje lozinki)
4. Za developere: biti pažljiv

**Beleške**:

Sve informacije za rešavanje izazova se nalaze u njegovom opisu i hintovima. Možemo pretpostaviti da će šifra nekako sadržati reč "kif". Guglovanjem izraza "93.83 billion trillion trillion centuries to brute force" nailazimo na stranicu za proveru jačine lozinki, a iz hinta o "One Important Final Note" zaključujemo da Će šifra verovatno sadržati u sebi tačke. Nakon nekoliko permutacija dolazimo do šifre "K1f....................."