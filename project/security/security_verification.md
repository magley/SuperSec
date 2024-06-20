У овом документу налази се анализа OWASP ASVS 4.0.3 документа примењен на Lesotho пројекат.

## V1 Architecture, Design and Threat Modeling

### V1.1 Secure Software Development Lifecycle

Није разматрано. Овај пројекат је сувише малих размера и разбијање на целине,
корисничке приче итд.; а ради дефинисања безбедносних апсеката није исплативо.

### V1.2 Authentication Architecture

| Број | Опис | CWE | Испуњеност |
| ---- | ---- | --- | ---------- |
| 1.2.1| Verify the use of unique or special low-privilege operating system accounts for all application components, services, and servers.| 250 | **Да**, ниједан сервис не захтева администраторске привилегије. |
| 1.2.2| Verify that communications between application components, including APIs, middleware and data layers, are authenticated. Components should have the least necessary privileges needed. | 306 | **Не**, интерне компоненте су подешене са подразумеваним привилегијама. Свако може приступити Lesotho севрису - **није имплементиран механизам API кључева**. |
| 1.2.3| Verify that the application uses a single vetted authentication mechanism that is known to be secure, can be extended to include strong authentication, and has sufficient logging and monitoring to detect account abuse or breaches. | 306 | **Не**, logging и monitoring нису имплементирани. |
| 1.2.4| Verify that all authentication pathways and identity management APIs implement consistent authentication security control strength, such that there are no weaker alternatives per the risk of the application.| 306 | **Да**, имплементиран је само један аутентификациони (_demo2_) и само један ауторизациони (_централизовани Lesotho auth_) систем у склопу пројекта.
---

### V1.3 Session Management Architecture

> This is a placeholder for future architectural requirements.

### V1.4 Access Control Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.4.1 | Verify that trusted enforcement points, such as access control gateways, servers, and serverless functions, enforce access controls. Never enforce access controls on the client. | 602 | **Да** |
| 1.4.2 | [DELETED, NOT ACTIONABLE] |
| 1.4.3 | [DELETED, DUPLICATE OF 4.1.3] |
| 1.4.4 | Verify the application uses a single and well-vetted access control mechanism for accessing protected data and resources. All requests must pass through this single mechanism to avoid copy and paste or insecure alternative paths. | 284 | **Да**, коришћен је Lesotho. Demo2 шаље ауторизационе захтеве серверу при покушају читања, писања и дељења документа. |
| 1.4.5 | Verify that attribute or feature-based access control is used whereby the code checks the user's authorization for a feature/data item rather than just their role. Permissions should still be allocated using roles. | 275 | **Не**, коришћен је наивни RBAC где се пермисије успостављају на нивоу рола |

### V1.5 Input and Output Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.5.1 | Verify that input and output requirements clearly define how to handle and process data based on type, content, and applicable laws, regulations, and other policy compliance. | 1029 | **?**
| 1.5.2 | Verify that serialization is not used when communicating with untrusted clients. If this is not possible, ensure that adequate integrity controls (and possibly encryption if sensitive data is sent) are enforced to prevent deserialization attacks including object injection. | 502 | **Да**, namespace и ACL се валидирају на пријему. |
| 1.5.3 | Verify that input validation is enforced on a trusted service layer. | 602 | **Да**, све валидације се извршавају на серверу. Потврда лозинке се ради на клијентском слоју, међутим то се не сматра безбедносним пропустом. |
| 1.5.4 | Verify that output encoding occurs close to or by the interpreter for which it is intended. | 116 | **?** |

### V1.6 Cryptographic Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.6.1 | Verify that there is an explicit policy for management of cryptographic keys and that a cryptographic key lifecycle follows a key management standard such as NIST SP 800-57. | 320 | **N/A** |
| 1.6.2 | Verify that consumers of cryptographic services protect key material and other secrets by using key vaults or API based alternatives. | 320 | **N/A** |
| 1.6.3 | Verify that all keys and passwords are replaceable and are part of a well-defined process to re-encrypt sensitive data. | 320 | **Не**, шифре нису заменљиве. |
| 1.6.4 | Verify that the architecture treats client-side secrets--such as symmetric keys, passwords, or API tokens--as insecure and never uses them to protect or access sensitive data. | 320 | **Да**, шифре се користе искључиво за аутентификацију. |

### V1.7 Errors, Logging and Auditing Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.7.1 | Verify that a common logging format and approach is used across the system. | 1009 | **Не**, логовање није имплементирано. |
| 1.7.2 | Verify that logs are securely transmitted to a preferably remote system for analysis, detection, alerting, and escalation. | | **Не**, логовање није имплементирано. |

### V1.8 Data Protection and Privacy Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.8.1 | Verify that all sensitive data is identified and classified into protection levels. | | **Не**, сви су у истом нивоу заштите (мрежном). |
| 1.8.2 | 1.8.2 Verify that all protection levels have an associated set of protection requirements, such as encryption requirements, integrity requirements, retention, privacy and other confidentiality requirements, and that these are applied in the architecture. | | **Не** |

### V1.9 Communications Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.9.1 | Verify the application encrypts communications between components, particularly when these components are in different containers, systems, sites, or cloud providers. | 319 | **Не**, HTTPs није имплементиран. |
| 1.9.2 | Verify that application components verify the authenticity of each side in a communication link to prevent person-in-the-middle attacks. For example, application components should validate TLS certificates and chains. | 295 | **Не**, HTTPs није имплементиран. |

### V1.10 Malicious Software Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.10.1 | Verify that a source code control system is in use, with procedures to ensure that check-ins are accompanied by issues or change tickets. The source code control system should have access control and identifiable users to allow traceability of any changes. | 284 | **Да**, коришћен је git, репозиторијум је на github серверу, и _issues_ су **делимично коришћени**. |

### V1.11 Business Logic Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.11.1 | Verify the definition and documentation of all application components in terms of the business or security functions they provide. | 1059 | **Да**, пословна логика функционалности, њихова документација и јединични тестови су усаглашћени.
| 1.11.2 | Verify that all high-value business logic flows, including authentication, session management and access control, do not share unsynchronized state. | 362 | **N/A**, над једном инстанцом је захтев испуњен. Скалирање сервиса није тестирано.
| 1.11.3 | Verify that all high-value business logic flows, including authentication, session management and access control are thread safe and resistant to time-of-check and time-of-use race conditions. | 367 | **? (Не)**, LevelDB и Consul-ов KV store су _thread-safe_, док Lesotho-ова кеш меморија није тестирана али вероватно није _thread-safe_. |

### V1.12 Secure File Upload Architecture

| Број   | Опис | CWE | Испуњеност |
| ------ | ---- | --- | ---------- |
| 1.12.1 | [DELETED, DUPLICATE OF 12.4.1] | |
| 1.12.2 | Verify that user-uploaded files - if required to be displayed or downloaded from the application - are served by either octet stream downloads, or from an unrelated domain, such as a cloud file storage bucket. Implement a suitable Content Security Policy (CSP) to reduce the risk from XSS vectors or other attacks from the uploaded file. | 646 | **Да**, фајлови се не приказују нити преузимају. |

### V1.13 API Architecture

> This is a placeholder for future architectural requirements.

### V1.14 Configuration Architecture

| Број   | Опис | CWE | Испуњеност |
| ------ | ---- | --- | ---------- |
| 1.14.1 | Verify the segregation of components of differing trust levels through well-defined security controls, firewall rules, API gateways, reverse proxies, cloud-based security groups, or similar mechanisms. | 923 | **N/A**, читав пројекат се извршава локално.
| 1.14.2 | Verify that binary signatures, trusted connections, and verified endpoints are used to deploy binaries to remote devices. | 494 | **N/A** |
| 1.14.3 | Verify that the build pipeline warns of out-of-date or insecure components and takes appropriate actions. | 1104 | **N/A** |
| 1.14.4 | Verify that the build pipeline contains a build step to automatically build and verify the secure deployment of the application, particularly if the application infrastructure is software defined, such as cloud environment build scripts. | | **Не**, рађено је ручно. |
| 1.14.5 | Verify that application deployments adequately sandbox, containerize and/or isolate at the network level to delay and deter attackers from attacking other applications, especially when they are performing sensitive or dangerous actions such as deserialization. | 265 | **Не** |
| 1.14.6 | Verify the application does not use unsupported, insecure, or deprecated client-side technologies such as NSAPI plugins, Flash, Shockwave, ActiveX, Silverlight, NACL, or client-side Java applets. | 477 | **Да** |

## V2 Authentication

### V2.1 Password Security

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.1.1 | Verify that user set passwords are at least 12 characters in length (after multiple spaces are combined) | 521 | 5.1.1.2 | **Да** |
| 2.1.2 | Verify that passwords of at least 64 characters are permitted, and that passwords of more than 128 characters are denied. | 521 | 5.1.1.2 | **Да** |
| 2.1.3 | Verify that password truncation is not performed. However, consecutive multiple spaces may be replaced by a single space. | 521 | 5.1.1.2 | **Да** |
| 2.1.4 | Verify that any printable Unicode character, including language neutral characters such as spaces and Emojis are permitted in passwords. | 521 | 5.1.1.2 | **Не**, коришћен је `input()`. |
| 2.1.5 | Verify users can change their password. | 620 | 5.1.1.2 | **Не** |
| 2.1.6 | Verify that password change functionality requires the user's current and new password. | 620 | 5.1.1.2 | **N/A** |
| 2.1.7 | Verify that passwords submitted during account registration, login, and password change are checked against a set of breached passwords either locally (such as the top 1,000 or 10,000 most common passwords which match the system's password policy) or using an external API. If using an API a zero knowledge proof or other mechanism should be used to ensure that the plain text password is not sent or used in verifying the breach status of the  password. If the password is breached, the application must require the user to set a new non-breached password. | 521 | 5.1.1.2 | **Не** |
| 2.1.8 | Verify that a password strength meter is provided to help users set  a stronger password. | 521 | 5.1.1.2 | **Не**, јер је коришћен CLI |
| 2.1.9 | Verify that there are no password composition rules limiting the  type of characters permitted. There should be no requirement for  upper or lower case or numbers or special characters. | 521 | 5.1.1.2 | **Да** |
| 2.1.10 | Verify that there are no periodic credential rotation or password  history requirements. | 263 | 5.1.1.2 | **Да** |
| 2.1.11 | Verify that "paste" functionality, browser password helpers, and  external password managers are permitted. | 521 | 5.1.1.2 | **Да** |
| 2.1.12 | Verify that the user can choose to either temporarily view the  entire masked password, or temporarily view the last typed character of the password on platforms that do not have this as built-in functionality. | 521 | 5.1.1.2 | **Не** |

### V2.2 General Authenticator Security

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.2.1 | Verify that anti-automation controls are effective at mitigating  breached credential testing, brute force, and account lockoutattacks. Such controls include blocking the most common  breached passwords, soft lockouts, rate limiting, CAPTCHA, ever  increasing delays between attempts, IP address restrictions, or  risk-based restrictions such as location, first login on a device,  recent attempts to unlock the account, or similar. Verify that no  more than 100 failed attempts per hour is possible on a single  account. | 307 | 5.2.2 / 5.1.1.2 / 5.1.4.2 / 5.1.5.2 | **Не** |
| 2.2.2 | Verify that the use of weak authenticators (such as SMS and email) is limited to secondary verification and transaction approval and not as a replacement for more secure authentication methods. Verify that stronger methods are offered before weak methods, users are aware of the risks, or that proper measures are in place to limit the risks of account compromise. | 304 | 5.2.10.0 | **Не** |
| 2.2.3 | Verify that secure notifications are sent to users after updates to authentication details, such as credential resets, email or address changes, logging in from unknown or risky locations. The use of push notifications - rather than SMS or email - is preferred, but in the absence of push notifications, SMS or email is acceptable as long as no sensitive information is disclosed in the notification. | 620 | | **Не** |
| 2.2.4 | Verify impersonation resistance against phishing, such as the use of multi-factor authentication, cryptographic devices with intent (such as connected keys with a push to authenticate), or at higher AAL levels, client-side certificates. | 308 | 5.2.5 | **Не**
| 2.2.5 | Verify that where a Credential Service Provider (CSP) and the  application verifying authentication are separated, mutually  authenticated TLS is in place between the two endpoints. | 319 | 5.2.6 | **Не** |
| 2.2.6 | Verify replay resistance through the mandated use of One-time Passwords (OTP) devices, cryptographic authenticators, or lookup codes. | 308 | 5.2.8 | **Не** |
| 2.2.7 | Verify intent to authenticate by requiring the entry of an OTP token or user-initiated action such as a button press on a FIDO hardware key. | 308 | 5.2.9 | **Не** |

### V2.3 Authenticator Lifecycle

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.3.1 | Verify system generated initial passwords or activation codes  SHOULD be securely randomly generated, SHOULD be at least 6 characters long, and MAY contain letters and numbers, and expire after a short period of time. These initial secrets must not be permitted to become the long term password. | 330 | 5.1.1.2 | **N/A** |
| 2.3.2 | Verify that enrollment and use of user-provided authentication  devices are supported, such as a U2F or FIDO tokens. | 308 | 6.1.3 | **Не**|
| 2.3.3 | Verify that renewal instructions are sent with sufficient time to  renew time bound authenticators. | 287 | 6.1.4 | **Не** |

### V2.4 Credential Storage

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.4.1 | Verify that passwords are stored in a form that is resistant to offline attacks. Passwords SHALL be salted and hashed using an approved one-way key derivation or password hashing function. Key derivation and password hashing functions take a password, a salt, and a cost factor as inputs when generating a password hash. | 916 | 5.1.1.2 | **Да**, коришћен је _bcrypt_. |
| 2.4.2 | Verify that the salt is at least 32 bits in length and be chosen arbitrarily to minimize salt value collisions among stored hashes. For each credential, a unique salt value and the resulting hash SHALL be stored. | 916 | 5.1.1.2 | **Да** |
| 2.4.3 | Verify that if PBKDF2 is used, the iteration count SHOULD be as large as verification server performance will allow, typically at least 100,000 iterations. | 916 | 5.1.1.2 | **N/A** |
| 2.4.4 | Verify that if bcrypt is used, the work factor SHOULD be as large as  verification server performance will allow, with a minimum of 10. | 916 | 5.1.1.2 | **Да**, коришћена је подразумевана вредност од 12 ([извор](https://pypi.org/project/bcrypt/#:~:text=Adjustable%20Work%20Factor)) |
| 2.4.5 | Verify that an additional iteration of a key derivation function is performed, using a salt value that is secret and known only to the verifier. Generate the salt value using an approved random bit generator [SP 800-90Ar1] and provide at least the minimum security strength specified in the latest revision of SP 800-131A. The secret salt value SHALL be stored separately from the hashed passwords (e.g., in a specialized device like a hardware security module). | 916 | 5.1.1.2 | **Не** |

### V2.5 Credential Recovery

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.5.1 | Verify that a system generated initial activation or recovery secret is not sent in clear text to the user. | 640 | 5.1.1.2 | **N/A** |
| 2.5.2 | Verify password hints or knowledge-based authentication (so-called "secret questions") are not present | 640 | 5.1.1.2 | **Да** |
| 2.5.3 | Verify password credential recovery does not reveal the current password in any way. | 640 | 5.1.1.2 | **N/A** |
| 2.5.4 | Verify shared or default accounts are not present (e.g. "root", "admin", or "sa"). | 16 | 5.1.1.2 | **Да** |
| 2.5.5 | Verify that if an authentication factor is changed or replaced, that the user is notified of this event. | 304 | 6.1.2.3 | **N/A** |
| 2.5.6 | Verify forgotten password, and other recovery paths use a secure recovery mechanism, such as time-based OTP (TOTP) or other soft token, mobile push, or another offline recovery mechanism. | 640 | 5.1.1.2 | **N/A** |
| 2.5.7 | Verify that if OTP or multi-factor authentication factors are lost, that evidence of identity proofing is performed at the same level as during enrollment. | 308 | 6.1.2.3 | **N/A** |

### V2.6 Look-up Secret Verifier

Није апликабилно за овај пројекат.

### V2.7 Out of Band Verifier

Није апликабилно за овај пројекат.

### V2.8 One Time Verifier

Није апликабилно за овај пројекат.

### V2.9 Cryptographic Verifier

Није апликабилно за овај пројекат.

### V2.10 Service Authentication

Није апликабилно за овај пројекат.

## V3 Session Management

### V3.1 Fundamental Session Management Security

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 3.1.1 | Verify the application never reveals session tokens in URL parameters. | 598 | | **Да** |

### V3.2 Session Binding

Није апликабилно за овај пројекат.

### V3.3 Session Termination

Није апликабилно за овај пројекат. Сесије су локалне и чувају се у радној меморији клијента (demo2).

### V3.4 Cookie-based Session Management

Није апликабилно за овај пројекат. Не користе се кукији.

### V3.5 Cookie-based Session Management

Није апликабилно за овај пројекат. Не користе се JWT/OAuth/SAML кључеви.

### V3.6 Federated Re-authentication

Није апликабилно за овај пројекат.

### V3.7 Defenses Against Session Management Exploits

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 3.7.1 | Verify the application ensures a full, valid login session or requires re-authentication or secondary verification before allowing any sensitive transactions or account modifications. | 306 | | **Не** |

## V4 Access Control

### V4.1 General Access Control Design

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 4.1.1 | Verify that the application enforces access control rules on a trusted service layer, especially if client-side access control is present and could be bypassed. | 602 | **Да** |
| 4.1.2 | Verify that all user and data attributes and policy information used by access controls cannot be manipulated by end users unless specifically authorized. | 639 | **Да** |
| 4.1.3 | Verify that the principle of least privilege exists - users should only be able to access functions, data files, URLs, controllers, services, and other resources, for which they possess specific authorization. This implies protection against spoofing and elevation of privilege.| 285 | **?** |
| 4.1.4 | [DELETED, DUPLICATE OF 4.1.3] | | |
| 4.1.5 | Verify that access controls fail securely including when an exception occurs. | 285 | **Не** | 

### V4.2 Operation Level Access Control

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 4.2.1 | Verify that sensitive data and APIs are protected against Insecure Direct Object Reference (IDOR) attacks targeting creation, reading, updating and deletion of records, such as creating or updating someone else's record, viewing everyone's records, or deleting all records. | 639 | **Да**, сваком захтеву претходи access control. |
| 4.2.2 | Verify that the application or framework enforces a strong anti-CSRF mechanism to protect authenticated functionality, and effective anti-automation or anti-CSRF protects unauthenticated functionalit | 352 | **Да**, CSRF није могућ. |

### V4.3 Other Access Control Considerations

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 4.3.1 | Verify administrative interfaces use appropriate multi-factor authentication to prevent unauthorized use. | 419 | **N/A** |
| 4.3.2 | Verify that directory browsing is disabled unless deliberately desired. Additionally, applications should not allow discovery or disclosure of file or directory metadata, such as Thumbs.db, .DS_Store, .git or .svn folders. | 548 | **N/A** |
| 4.3.3 | Verify the application has additional authorization (such as step up or adaptive authentication) for lower value systems, and / or segregation of duties for high value applications to enforce anti-fraud controls as per the risk of application and past fraud. | 732 | **N/A** |

## V5 Validation, Sanitization and Encoding

### V5.1 Input Validation

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.1.1 | Verify that the application has defenses against HTTP parameter pollution  attacks, particularly if the application framework makes no distinction about  the source of request parameters (GET, POST, cookies, headers, or  environment variables). | 235 | **Да** |
| 5.1.2 | Verify that frameworks protect against mass parameter assignment attacks,  or that the application has countermeasures to protect against unsafe  parameter assignment, such as marking fields private or similar. | 915 | **Да** |
| 5.1.3 | Verify that all input (HTML form fields, REST requests, URL parameters, HTTP  headers, cookies, batch files, RSS feeds, etc) is validated using positive  validation (allow lists). | 20 | **Не** |
| 5.1.4 | Verify that structured data is strongly typed and validated against a defined  schema including allowed characters, length and pattern (e.g. credit card  numbers, e-mail addresses, telephone numbers, or validating that two related  fields are reasonable, such as checking that suburb and zip/postcode match). | 20 | **Не** (парцијална валидација постоји) |
| 5.1.5 | Verify that URL redirects and forwards only allow destinations which appear  on an allow list, or show a warning when redirecting to potentially untrusted  content. | 601 | **N/A** |

### V5.2 Sanitization and Sandboxing

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.2.1 | Verify that all untrusted HTML input from WYSIWYG editors or similar is  properly sanitized with an HTML sanitizer library or framework feature. | 116 | **N/A** |
| 5.2.2 | Verify that unstructured data is sanitized to enforce safety measures such as  allowed characters and length. | 138 | **Не** |
| 5.2.3 | Verify that the application sanitizes user input before passing to mail systems  to protect against SMTP or IMAP injection. | 147 | **N/A** |
| 5.2.4 | Verify that the application avoids the use of eval() or other dynamic code  execution features. Where there is no alternative, any user input being  included must be sanitized or sandboxed before being executed. | 95 | **Да** |
| 5.2.5 | Verify that the application protects against template injection attacks by  ensuring that any user input being included is sanitized or sandboxed. | 94 | **N/A** |
| 5.2.6 | Verify that the application protects against SSRF attacks, by validating or  sanitizing untrusted data or HTTP file metadata, such as filenames and URL  input fields, and uses allow lists of protocols, domains, paths and ports. | 918 | **N/A** |
| 5.2.7 | Verify that the application sanitizes, disables, or sandboxes user-supplied  Scalable Vector Graphics (SVG) scriptable content, especially as they relate to  XSS resulting from inline scripts, and foreignObject. | 159 | **N/A** |
| 5.2.8 | Verify that the application sanitizes, disables, or sandboxes user-supplied  scriptable or expression template language content, such as Markdown, CSS  or XSL stylesheets, BBCode, or similar. | 94 | **N/A** |

### V5.3 Output Encoding and Injection Prevention

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.3.1 | Verify that output encoding is relevant for the interpreter and context  required. For example, use encoders specifically for HTML values, HTML  attributes, JavaScript, URL parameters, HTTP headers, SMTP, and others as  the context requires, especially from untrusted inputs (e.g. names with  Unicode or apostrophes, such as ねこ or O'Hara). | 116 | **Не** |
| 5.3.2 | Verify that output encoding preserves the user's chosen character set and  locale, such that any Unicode character point is valid and safely handled. | 176 | **Не** |
| 5.3.3 | Verify that context-aware, preferably automated - or at worst, manual - output escaping protects against reflected, stored, and DOM based XSS. | 79 | **Не** |
| 5.3.4 | Verify that data selection or database queries (e.g. SQL, HQL, ORM, NoSQL)  use parameterized queries, ORMs, entity frameworks, or are otherwise  protected from database injection attacks. | 89 | **N/A** |
| 5.3.5 | Verify that where parameterized or safer mechanisms are not present,  context-specific output encoding is used to protect against injection attacks,  such as the use of SQL escaping to protect against SQL injection. | 89 | **N/A** |
| 5.3.6 | Verify that the application protects against JSON injection attacks, JSON eval  attacks, and JavaScript expression evaluation. | 830 | **?** |
| 5.3.7 | Verify that the application protects against LDAP injection vulnerabilities, or  that specific security controls to prevent LDAP injection have been  implemented. | 90 | **N/A** |
| 5.3.8 | Verify that the application protects against OS command injection and that  operating system calls use parameterized OS queries or use contextual  command line output encoding. | 78 | **N/A** |
| 5.3.9 | Verify that the application protects against Local File Inclusion (LFI) or  Remote File Inclusion (RFI) attacks. | 829 | **N/A** |
| 5.3.10 | Verify that the application protects against XPath injection or XML injection  attacks.  | 643 | **N/A** |

### V5.4 Memory, String, and Unmanaged Code

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.4.1 | Verify that the application uses memory-safe string, safer memory copy and  pointer arithmetic to detect or prevent stack, buffer, or heap overflows. | 120 | **?** |
| 5.4.2 | Verify that format strings do not take potentially hostile input, and are  constant. | 134 | **Не** (парцијално задовољено) |
| 5.4.3 | Verify that sign, range, and input validation techniques are used to prevent  integer overflows. | 190 | **N/A** |

### V5.5 Deserialization Prevention

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.5.1 | Verify that serialized objects use integrity checks or are encrypted to prevent  hostile object creation or data tampering. | 502 | **Да** |
| 5.5.2 | Verify that the application correctly restricts XML parsers to only use the most  restrictive configuration possible and to ensure that unsafe features such as  resolving external entities are disabled to prevent XML eXternal Entity (XXE)  attacks. | 611 | **N/A** |
| 5.5.3 | Verify that deserialization of untrusted data is avoided or is protected in both  custom code and third-party libraries (such as JSON, XML and YAML parsers). | 502 | **Не** |
| 5.5.4 | Verify that when parsing JSON in browsers or JavaScript-based backends,  JSON.parse is used to parse the JSON document. Do not use eval() to parse  JSON. | 95 | **Да** |

## V6 Stored Cryptography

Ово поглавље није разматрано за овај пројекат.

## V7 Error Handling and Logging

### V7.1 Log Content

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.1.1 | Verify that the application does not log credentials or payment details.  Session tokens should only be stored in logs in an irreversible, hashed form. | 532 | **Да** |
| 7.1.2 | Verify that the application does not log other sensitive data as defined under  local privacy laws or relevant security policy. | 532 | **Да** |
| 7.1.3 | Verify that the application logs security relevant events including successful  and failed authentication events, access control failures, deserialization  failures and input validation failures. | 778 | **Не** |
| 7.1.4 | Verify that each log event includes necessary information that would allow for  a detailed investigation of the timeline when an event happens. | 778 | **Не** |

### V7.2 Log Processing

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.2.1 | Verify that all authentication decisions are logged, without storing sensitive  session tokens or passwords. This should include requests with relevant  metadata needed for security investigations. | 778 | **Не** |
| 7.2.2 | Verify that all access control decisions can be logged and all failed decisions  are logged. This should include requests with relevant metadata needed for  security investigations. | 285 | **Не** |

### V7.3 Log Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.3.1 | Verify that all logging components appropriately encode data to prevent log  injection. | 117 | **N/A** |
| 7.3.2 | [DELETED, DUPLICATE OF 7.3.1] |  | |
| 7.3.3 | Verify that security logs are protected from unauthorized access and  modification. | 200 | **N/A** |
| 7.3.4 | Verify that time sources are synchronized to the correct time and time zone.  Strongly consider logging only in UTC if systems are global to assist with post-incident forensic analysis. |  | **Не** |

### V7.4 Error Handling

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.4.1 | Verify that a generic message is shown when an unexpected or security  sensitive error occurs, potentially with a unique ID which support personnel  can use to investigate.  | 210 | **Не**, порука је дескриптивна |
| 7.4.2 | Verify that exception handling (or a functional equivalent) is used across the  codebase to account for expected and unexpected error conditions. | 533 | **Да** |
| 7.4.3 | Verify that a "last resort" error handler is defined which will catch all  unhandled exceptions.  | 431 | **?**, није познато за сам Lesotho сервис где је коришћен Go. |

## V8 Data Protection

### V8.1 General Data Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 8.1.1 | Verify the application protects sensitive data from being cached in server  components such as load balancers and application caches. | 524 | **Да**, кеширање је ручно и једино се кешира namespace који нема осетљиве податке. |
| 8.1.2 | Verify that all cached or temporary copies of sensitive data stored on the  server are protected from unauthorized access or purged/invalidated after  the authorized user accesses the sensitive data. | 524 | **N/A**, |
| 8.1.3 | Verify the application minimizes the number of parameters in a request, such  as hidden fields, Ajax variables, cookies and header values. | 233 | **Да** |
| 8.1.4 | Verify the application can detect and alert on abnormal numbers of requests,  such as by IP, user, total per hour or day, or whatever makes sense for the  application. | 770 | **Не** |
| 8.1.5 | Verify that regular backups of important data are performed and that test  restoration of data is performed. | 19 | **Не** |
| 8.1.6 | Verify that backups are stored securely to prevent data from being stolen or  corrupted. | 19 | **N/A** |

### V8.2 Client-side Data Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 8.2.1 | Verify the application sets sufficient anti-caching headers so that sensitive  data is not cached in modern browsers. | 525 | **N/A** |
| 8.2.2 | Verify that data stored in browser storage (such as localStorage,  sessionStorage, IndexedDB, or cookies) does not contain sensitive data. | 922 | **N/A** |
| 8.2.3 | Verify that authenticated data is cleared from client storage, such as the  browser DOM, after the client or session is terminated. | 922 | **N/A** |

### V8.3 Sensitive Private Data

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 8.3.1 | Verify that sensitive data is sent to the server in the HTTP message body or  headers, and that query string parameters from any HTTP verb do not contain  sensitive data. | 319 | **Да** |
| 8.3.2 | Verify that users have a method to remove or export their data on demand.  | 212 | **Не** |
| 8.3.3 | Verify that users are provided clear language regarding collection and use of  supplied personal information and that users have provided opt-in consent  for the use of that data before it is used in any way. | 285 | **Не** |
| 8.3.4 | Verify that all sensitive data created and processed by the application has  been identified, and ensure that a policy is in place on how to deal with  sensitive data. | 200 | **Не** |
| 8.3.5 | Verify accessing sensitive data is audited (without logging the sensitive data  itself), if the data is collected under relevant data protection directives or  where logging of access is required. | 532 | **N/A** |
| 8.3.6 | Verify that sensitive information contained in memory is overwritten as soon  as it is no longer required to mitigate memory dumping attacks, using zeroes  or random data. | 226 | **N/A** |
| 8.3.7 | Verify that sensitive or private information that is required to be encrypted, is  encrypted using approved algorithms that provide both confidentiality and  integrity.  | 327 | **Да** |
| 8.3.8 | Verify that sensitive personal information is subject to data retention  classification, such that old or out of date data is deleted automatically, on a  schedule, or as the situation requires. | 285 | **N/A** |

## V9 Communication