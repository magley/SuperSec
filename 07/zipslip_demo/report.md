Kreirana je jednostavna _Python Flask_ aplikacija sa dva endpointa u `app.py` fajlu:
- `/`: Sadrži formu za upload tar arhiva
- `/upload`: Raspakuje i prikazuje sadržaj arhive

Aplikacija je pokrenuta u _debug_ režimu kako bi se omogućila automatska detekcija izmena fajlova.

Kod za raspakivanje arhive je preuzet sa materijala sa vežbi `Slippy-CVE-2007-4559.pdf`. Razlike su promenjena putanja raspakivanja, i linija 
```python
os.rename(os.path.join(tmp, name), filename)
```
zamenjena je sa
```python
os.replace(os.path.join(tmp, name), filename)
```
Razlog tome je što _Windows_ sa prvom linijom ne dozvoljava kreiranje fajla ako taj fajl već postoji, što predstavlja problem ako želimo da demonstriramo ranjivost.

Ideja iza ranjivosti je da arhive mogu sadržati fajlove čije putanje mogu sadržati _directory traversal_ takav da se fajl raspakuje van trenutnog direktorijuma. Dakle, prilikom raspakovanja arhive možemo zameniti neki bitan fajl. Konkretno, mi ćemo zameniti `app.py` koji sadrži našu celu _Flask_ aplikaciju. Pošto aplikaciju pokrećemo u _debug_ režimu ova promena će se automatski prepoznati i naš novi `app.py` će se učitati, bez potrebe da se server restartuje.

Na osnovu skripte za pravljenje maliciozne arhive u `Slippy-CVE-2007-4559.pdf` napravili smo `make_harmful_tar.py`. Ova skripta pravi kompresovanu arhivu sa novim `app.py` fajlom koji će se raspakovati u direktorijum iznad trenutnog (gde se i nalazi tekući `app.py`). *Bitno* je da vreme modifikacije fajla (u skripti `info.mtime`) bude skorije od tekućeg `app.py` kako bi _Flask_ prepoznao izmenu.

Koraci za demonstraciju ranjivosti:

1. Pokrenuti aplikaciju sa `flask --app app run --debug`
2. Otići na `localhost:5000`
3. _Submit_ `harmless.tar.gz` za primer namenjenog ponašanja aplikacije
4. Generisati `harmful.tar.gz` pomoću `make_harmful_tar.py` skripte
5. _Submit_ `harmful.tar.gz`, osvežiti `localhost:5000`