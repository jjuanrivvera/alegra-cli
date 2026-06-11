---
title: Autocompletado del shell
---

# Autocompletado del shell

`alegra` trae un autocompletado completo para **bash**, **zsh**, **fish** y
**PowerShell**. Más allá de completar nombres de comandos y flags, se comunica
con la API para completar *tus datos reales*: presiona <kbd>Tab</kbd> después de
`invoices get` y obtienes IDs reales de facturas, cada uno anotado con su estado.

!!! tip "¿Ya está instalado?"
    Los archivos comprimidos de release, los paquetes `.deb`/`.rpm`/`.apk` y las
    instalaciones de Homebrew/Scoop incluyen los scripts de autocompletado. Si
    instalaste de esa forma, puede que el autocompletado ya funcione: abre un
    shell nuevo y prueba `alegra <Tab>`. Las secciones de abajo son para builds
    con `go install` / desde el código fuente, o para (re)instalarlo a mano.

## Qué se completa

| Lo que escribes | <kbd>Tab</kbd> ofrece |
| --- | --- |
| `alegra <Tab>` | cada recurso y comando (`invoices`, `contacts`, `auth`, …) |
| `alegra invoices <Tab>` | `list`, `get`, `create`, `update`, `delete`, `export`, `void`, … |
| `alegra invoices get <Tab>` | **IDs de facturas en vivo**, etiquetados con nombre/número/estado |
| `alegra invoices delete 10<Tab>` | IDs en vivo que empiezan con `10` |
| `alegra invoices void <Tab>` | IDs en vivo (las acciones personalizadas también completan IDs) |
| `alegra invoices list --status <Tab>` | `open`, `closed`, `draft`, `void` |
| `alegra invoices list --order-field <Tab>` | los campos ordenables del recurso |
| `alegra contacts get --columns <Tab>` | las columnas de ese recurso (reconoce las comas) |
| `alegra -o <Tab>` | `table`, `json`, `yaml`, `csv` |
| `alegra --profile <Tab>` | tus perfiles configurados, etiquetados con su email |

Los completadores que conocen tus datos (IDs, columnas, perfiles, enums) hacen
que la potencia de la CLI sea fácil de descubrir sin memorizar flags ni leer la
documentación.

!!! note "Cómo se comporta el autocompletado en vivo"
    El autocompletado de IDs hace una sola petición corta (de una página) a la
    API usando tu perfil activo. Es **de mejor esfuerzo y nunca bloquea tu
    shell**: sin credenciales, sin conexión o en una red lenta simplemente no
    devuelve sugerencias (un timeout de 4 segundos limita la espera) y puedes
    escribir el ID a mano. Respeta `--profile`, así que
    `alegra --profile sandbox invoices get <Tab>` completa contra tu cuenta de
    sandbox.

## Instalación

=== "bash"

    El autocompletado necesita el paquete `bash-completion`.

    ```bash
    # Current shell only (try it now):
    source <(alegra completion bash)

    # Persist — Linux:
    alegra completion bash | sudo tee /etc/bash_completion.d/alegra >/dev/null

    # Persist — macOS (Homebrew bash-completion@2):
    alegra completion bash > "$(brew --prefix)/etc/bash_completion.d/alegra"
    ```

=== "zsh"

    Asegúrate de habilitar el autocompletado una vez (en `~/.zshrc`):

    ```bash
    echo "autoload -U compinit; compinit" >> ~/.zshrc
    ```

    Luego instala el script en tu `$fpath`:

    ```bash
    # A common location; create it if needed and ensure it's on your fpath.
    alegra completion zsh > "${fpath[1]}/_alegra"
    ```

    Abre un shell nuevo. Con Oh My Zsh, escribe en
    `~/.oh-my-zsh/completions/_alegra` en su lugar.

=== "fish"

    ```bash
    # Current shell only:
    alegra completion fish | source

    # Persist:
    alegra completion fish > ~/.config/fish/completions/alegra.fish
    ```

=== "PowerShell"

    ```powershell
    # Current shell only:
    alegra completion powershell | Out-String | Invoke-Expression

    # Persist by appending to your profile:
    alegra completion powershell >> $PROFILE
    ```

Ejecuta `alegra completion <shell> --help` para ver las instrucciones
autoritativas y específicas de cada shell.

## Verificar

Abre un shell **nuevo** y prueba:

```bash
alegra invoices <Tab><Tab>     # subcommands
alegra -o <Tab>                # table / json / yaml / csv
alegra invoices get <Tab>      # live IDs (requires a logged-in profile)
```

Si no pasa nada, confirma que el autocompletado esté habilitado para tu shell
(bash necesita `bash-completion`; zsh necesita `compinit`) y que hayas iniciado
una sesión nueva.

## Solución de problemas

- **No hay sugerencias después de `get`/`delete`**: estos son autocompletados en
  vivo. Verifica que tengas un perfil funcionando (`alegra auth status`); sin
  credenciales no devuelven nada por diseño.
- **zsh: `command not found: compdef`**: agrega `autoload -U compinit; compinit`
  a `~/.zshrc` y reinicia el shell.
- **Autocompletados desactualizados tras una actualización**: regenera el script
  (los comandos de arriba) y abre un shell nuevo.
