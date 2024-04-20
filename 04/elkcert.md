# Pronađeni defekti i preporuke za poboljšanje

Projekat: sssscs frontend (videti sssscs.md)

## (1) Nedovoljno korišćenje `const` promenljivih

Preferirajući `const` promenljive eksplicitno dajemo do znanja programeru da promenljiva nije namenjena za menjanje. Ovo je korisno radi čitljivosti i razumevanja koda jer jasnije izražava kako se promenljiva koristi, a i generalno je bezbednija praksa (opt-in umesto opt-out pristup, ne omogućavamo nepotrebnu mutaciju podataka).

Ovaj problem se prostire kroz ceo projekat. Rešenje je zamena `let` ključnom reči sa `const` na svim označenim mestima (linter ovo može uraditi automatski)

## (2) 'React' must be in scope when using JSX

Nije tehnički defekt u kodu, već zabuna u linteru. [Istraživanjem](https://kinsta.com/knowledgebase/react-must-be-in-scope-when-using-jsx/) saznajemo da do određene verzije React-a je bilo potrebno importovati sledeće:

```typescript
import React  from 'react';
```

Međutim, od verzije 17 ovaj import se obavi implicitno (mi koristimo verziju 18.2). Rešenje je da u eslint config-u ili isključimo ovu proveru, ili da obavestimo config koju verziju React-a koristimo (linter se žali da mu nije prosleđena verzija). Međutim, nijedan isprobani način nije uspeo, verovatno zbog specifičnog načina na koji je projekat podešen.

## (3) Definisane neiskorišćene promenljive

Na iznenađujuće puno mesta smo definisali neiskorišćene promenljive. Ako želimo da kod prati bilo koju kodnu konvenciju, ovo bi trebalo izbegavati.


## (4) Upotreba `any` tipa

U fajlu http/HttpService.ts:

```typescript
axiosInstance.interceptors.response.use(
    (response: AxiosResponse<any, any>) => {
        return response;
    },
    (error: any) => {
        if (error.response.status === 401) {
            logoutAndMoveToLoginPage();
        }
        return Promise.reject(error);
    }
)
```

Jedna od prednosti (i razloga) upotrebe Typescript-a je upravo povećana tipska bezbednost koju jezik pruža. Većina slučajeva upotrebe `any` tipa (uključujući i ovde) je lenjost.

## (5) Upotreba neprimitivnog `Boolean` tipa

U fajlu certs/CertVerifyFile.tsx:
```typescript
let [ isValid, setIsValid ] = useState<Boolean | null>(null);
```

Upotreba neprimitivnih tipova [nije preporučena](https://www.typescriptlang.org/docs/handbook/declaration-files/do-s-and-don-ts.html#general-types) u većini slučajeva.

## (6) Globalno stanje se prosleđuje kao props

Fajl App.tsx:

```typescript
// TODO: Explore the Context API. 
// We can share GlobalState without explicitly passing props.

export interface GlobalState {
    updateIsLoggedIn: () => void;
    isLoggedIn: boolean;
    role: string;
};

...

let [state, setState] = React.useState<GlobalState>({
    isLoggedIn: AuthService.isLoggedIn(),
    role: AuthService.getRole(),
    updateIsLoggedIn: () => {
        let newState = {...state};
        newState.isLoggedIn = AuthService.isLoggedIn();
        newState.role = AuthService.getRole();

        setState(newState);
    }
});
```

Umesto da globalno stanje prosleđujemo svim zainteresovanim komponentama (što može biti veoma mukotrpno), bolje rešenje bi bilo koristiti [Context API](https://react.dev/learn/passing-data-deeply-with-context).

# Statistika

**Provedeno vreme na analizi**: 2h

**Broj defekata**: 6

**Korišćeni alati**:

- [ESLint](https://eslint.org/)
