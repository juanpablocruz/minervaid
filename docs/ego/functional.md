# Documento funcional: Proyecto `ego`

## Nombre del proyecto

**ego** – Identidad digital autocontenida (CLI)

---

## Objetivo

Construir una herramienta CLI en Go que permita a un usuario generar, gestionar y utilizar identidades digitales autocontenidas siguiendo principios de SSI (Self-Sovereign Identity). Cada identidad es totalmente local, privada, portátil y puede operar tanto en contextos cerrados como abiertos (por ejemplo, usando `did:key`, `did:web`, etc.).

---

## Conceptos clave

### Vault

Un contenedor local (carpeta) que actúa como "wallet" de identidades. Contiene:

* Identidades independientes.
* Configuración general.

### Identidad

Una carpeta dentro del vault que representa:

* Un DID (con su tipo: key, web, etc.).
* Claves privadas (almacenadas como `keystore.json`).
* Atributos opcionales (`profile.json`).
* (Opcional) credenciales verificables (`vcs/`).

### DID

Identificador descentralizado generado localmente. Puede ser:

* `did:key`
* `did:web`
* (futuro: `did:ion`, `did:ebsi`, etc.)

### Atributos

Información asociada a la identidad, no necesariamente verificable (nombre, email, alias, etc.).

### Firmas y autenticación

Capacidad de firmar archivos o responder a desafíos (challenges) con pruebas criptográficas que usan el DID y la clave privada asociada.

### Motor SSI

La funcionalidad técnica subyacente de SSI (generación de DIDs, creación y firma de credenciales, etc.) es realizada mediante el uso de **MinervaID**, que actúa como motor SSI del sistema. `ego` sirve como capa de experiencia de usuario sobre esta base técnica.

---

## Estructura de archivos esperada

```plaintext
vault/
├── local-ego/                     # Identidad local
│   ├── did.txt                    # DID generado
│   ├── keystore.json              # Clave privada cifrada
│   ├── profile.json               # Atributos
│   └── vcs/                       # Credenciales (futuro)
├── ego-web/                      # Identidad migrada (did:web)
│   └── ...
├── config.json                   # Config global (identidad activa)
└── .current                      # Nombre de la identidad activa
```

---

## Comandos principales

### Inicialización y gestión de identidades

```bash
ego init --name <nombre> [--web dominio.com]     # Crea una nueva identidad
ego use <nombre>                                 # Selecciona identidad activa
ego current                                      # Muestra identidad activa
ego list                                         # Lista identidades
ego migrate --to web --domain dominio.com        # Clona identidad como did:web
```

### Perfil y atributos

```bash
ego set <clave> <valor>                          # Modifica atributo
ego show                                         # Muestra DID y atributos
```

### Firma

```bash
ego sign archivo.txt                             # Firma un archivo
ego sign archivo.txt --output firma.json         # Firma estructurada
ego sign archivo.txt --with rol=editor           # Adjunta atributos firmados
```

### Autenticación

```bash
ego auth respond --file challenge.json           # Firma un challenge SSI
```

### Remotos

```bash
ego remote add web https://miweb.com/.well-known/did.json
ego remote export web --out ./did.json           # Exporta documento DID
```

---

## Seguridad

* El `keystore.json` está cifrado y requiere desbloqueo manual con passphrase.
* Las claves privadas se usan solo en memoria durante operaciones autorizadas.
* Futuro: integración con FIDO2, huella o 2FA para desbloqueo.

---

## Integración futura: `ego-agent`

* Daemon que mantiene claves desbloqueadas en memoria.
* Expone socket o gRPC para firma bajo demanda.
* Usado por herramientas externas (SSH, CI, scripts).

---

## Casos de uso

* Gestión de identidades personales o seudónimas.
* Firma de documentos, publicaciones, artefactos.
* Autenticación SSI en apps web (challenge-response).
* Identidad portátil con soporte para múltiples métodos DID.

---

## Estado actual

* Interfaz definida.
* Estructura de carpetas establecida.
* Flujo funcional descrito.
* A la espera de implementación progresiva.

---

## Futuras extensiones

* Soporte para VCs y presentaciones selectivas.
* Interoperabilidad DIDComm / OpenID4VC.
* Firma colectiva (multi-identidad).
* Interfaz web o TUI para facilitar interacción.

---

## Autores

Proyecto `ego`, parte del stack de MinervaID.
