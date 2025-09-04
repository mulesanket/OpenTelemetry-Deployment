## Industry standard pipeline flow (Java):

Checkout

Pre-build scans (secrets, lint, basic SAST)

Compile (Gradle build)

Run tests + coverage

SonarQube analysis (using compiled code + coverage reports)

Package (JAR or installDist)

Docker build (runtime-only image)

Image scan (Trivy/Grype)