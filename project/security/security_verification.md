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

### V2 Authentication

