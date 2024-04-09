## 1)

**Naziv izazova u aplikaciji**: **Login Jim**

**Klasa napada**: SQL injekcija

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

**Naziv izazova u aplikaciji**: ****

**Klasa napada**:

**Uticaj**:

**Ranjivosti**:

**Kontramere**:

**Beleške**:

## 5)

**Naziv izazova u aplikaciji**: ****

**Klasa napada**:

**Uticaj**:

**Ranjivosti**:

**Kontramere**:

**Beleške**:
