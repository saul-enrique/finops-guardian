# 🛡️ FinOps Guardian 🛡️

Este proyecto es un bot de CI/CD para GitHub Actions que analiza automáticamente el impacto en los costos de los cambios de infraestructura en un Pull Request. Utiliza Terraform e Infracost para proporcionar una estimación clara de los costos antes de que los cambios se fusionen, promoviendo una cultura de conciencia de costos (FinOps) en los equipos de desarrollo.

## ¿Qué Problema Resuelve?

En los flujos de trabajo de DevOps modernos, los desarrolladores pueden realizar cambios en la infraestructura como código (IaC) sin ser conscientes del impacto económico que estos cambios tendrán en la factura de la nube. "FinOps Guardian" aborda este problema haciendo que el costo sea un "riesgo duro" y visible directamente en el flujo de trabajo del desarrollador.

## ¿Cómo Funciona?

Cuando se abre un Pull Request que modifica archivos de Terraform (`.tf`), se activa una GitHub Action que:
1.  Realiza un checkout tanto de la rama base como de la rama del PR.
2.  Utiliza la **acción oficial de Infracost** para comparar las dos versiones de la infraestructura.
3.  Calcula la diferencia de costos mensuales.
4.  Publica un comentario claro y conciso en el Pull Request con un resumen de los costos, permitiendo a los revisores tomar decisiones informadas.

![Ejemplo de Comentario](URL_A_UNA_CAPTURA_DE_PANTALLA_DEL_COMENTARIO_FINAL)

## Tecnologías Utilizadas

*   **Terraform:** Para definir la Infraestructura como Código.
*   **Infracost:** Para la estimación de costos de la nube.
*   **GitHub Actions:** Como plataforma de CI/CD para la automatización.
*   **YAML:** Para definir el workflow de la GitHub Action.
*   **(Opcional: Go):** Inicialmente, se desarrolló una CLI personalizada en Go para orquestar el proceso, lo que demostró una comprensión profunda de la automatización antes de refactorizar para usar la acción oficial por su robustez.

## El Viaje y las Lecciones Aprendidas

Este proyecto fue un viaje de depuración realista a través de múltiples capas de tecnología. Los desafíos incluyeron:
*   Manejar las complejidades de los entornos de `git` en CI/CD, como el estado "detached HEAD".
*   Asegurar la captura limpia de la salida de herramientas de línea de comandos.
*   Finalmente, comprender y aplicar la solución más robusta y mantenible utilizando las herramientas oficiales del ecosistema (la acción de Infracost), lo que refleja un principio clave de la ingeniería de software: **no reinventar la rueda cuando existe una solución profesional y probada.**

---
*Proyecto guiado por el asistente de IA.*
