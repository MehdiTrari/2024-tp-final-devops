# Rapport de Déploiement avec Docker, Docker Compose et Intégration Continue

## Introduction
Dans ce projet, j'ai configuré et déployé une application multi-services composée de :
- Une API (`vote-api`),
- Une interface utilisateur (`web-client`),
- Une documentation statique (`docs`),
- Une base de données PostgreSQL (`db`).

J'ai utilisé Docker pour conteneuriser chaque service, Docker Compose pour orchestrer leur déploiement, et une pipeline d'intégration continue (CI) pour automatiser les tests et la validation du code.

Enfin, j'ai tenté de déployer le projet sur Google Cloud Platform (GCP), mais je n'ai pas pu finaliser cette étape en raison de contraintes de temps et de difficultés techniques rencontrées dans la configuration.

---

## Méthode de Développement Agile et Gestion des Branches

Je travaille avec la méthode agile, ce qui implique une organisation stricte des branches pour faciliter le développement collaboratif :
- **Branche `main`** : Contient la version stable de l'application, prête à être déployée en production.
- **Branche `develop`** : Utilisée pour regrouper toutes les fonctionnalités en cours de développement avant leur validation finale.
- **Branches `feature/*`** : Créées à partir de `develop` pour chaque nouvelle fonctionnalité ou tâche. Une fois le travail terminé, elles sont fusionnées dans `develop`.
- **Hotfixes** : Les corrections urgentes (`hotfix`) sont créées directement depuis `main` pour résoudre des bugs critiques.

### Validation des Pull Requests
J'ai configuré GitHub pour :
- Exiger une validation automatique par les tests CI avant la fusion des PR.
- Empêcher la fusion si des tests échouent ou si la branche n'a pas été approuvée par un relecteur.

Cette configuration assure une gestion fluide et cohérente des versions tout en maintenant la stabilité de la branche principale.

---

## Déploiement des Images Docker

### Construction et Push des Images Docker
J'ai conteneurisé chaque service en utilisant Docker et j'ai publié les images sur Docker Hub pour faciliter leur déploiement. Voici les étapes effectuées :

1. **Construction des Images** :
   J'ai utilisé les commandes suivantes pour construire les images Docker :
   ```bash
   docker build -t mehditrr/vote-api:latest ./vote-api
   docker build -t mehditrr/web-client:latest ./web-client
   docker build -t mehditrr/docs:latest ./docs
   ```

2. **Publication des Images sur Docker Hub** :
   Les images construites ont été poussées sur Docker Hub à l'aide de ces commandes :
   ```bash
   docker push mehditrr/vote-api:latest
   docker push mehditrr/web-client:latest
   docker push mehditrr/docs:latest
   ```

3. **Liens des Images Docker** :
   - [vote-api](https://hub.docker.com/repository/docker/mehditrr/vote-api/general)
   - [web-client](https://hub.docker.com/repository/docker/mehditrr/web-client/general)
   - [docs](https://hub.docker.com/repository/docker/mehditrr/docs/general)

### Intégration dans la Pipeline CI
Pour automatiser le processus de build et de push des images Docker, j'ai ajouté un job `push-docker-images` dans la pipeline CI. Ce job est exécuté après la validation des tests et le build des artefacts :

```yaml
push-docker-images:
  runs-on: ubuntu-latest
  needs: build-artifacts
  steps:
    - name: Check out code
      uses: actions/checkout@v3

    - name: Log in to Docker Hub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}

    - name: Push vote-api image
      run: |
        docker build -t mehditrr/vote-api:latest ./vote-api
        docker push mehditrr/vote-api:latest

    - name: Push web-client image
      run: |
        docker build -t mehditrr/web-client:latest ./web-client
        docker push mehditrr/web-client:latest

    - name: Push docs image
      run: |
        docker build -t mehditrr/docs:latest ./docs
        docker push mehditrr/docs:latest
```

---

## Tentative de Déploiement sur GCP

J'ai configuré Google Cloud Platform pour déployer les services via **Cloud Run** et **Artifact Registry** :
1. J'ai créé un dépôt Docker sur Artifact Registry pour héberger mes images Docker.
2. J'ai tenté de déployer les images via Cloud Run, en passant les variables d'environnement nécessaires, notamment `PG_URL` pour connecter l'API à la base de données PostgreSQL.

### Problèmes rencontrés :
- **Connexion à PostgreSQL** : Le conteneur `vote-api` n'a pas pu établir une connexion avec la base de données, probablement en raison d'une configuration incorrecte ou de restrictions réseau.
- **Timeout des services** : Cloud Run a signalé que le conteneur ne répondait pas sur le port attendu dans le délai imparti.

En raison du manque de temps, je n'ai pas pu résoudre ces problèmes, mais j'ai identifié des pistes pour les prochaines étapes :
- Utiliser Google Cloud SQL pour gérer la base de données.
- Vérifier la configuration réseau pour autoriser Cloud Run à accéder à la base de données.

---

## Documentation pour Utilisateurs et Contributeurs

### Documentation Utilisateur
Pour exécuter l'application localement :
1. Installez Docker et Docker Compose.
2. Clonez le dépôt :
   ```bash
   git clone <repository_url>
   cd <repository_folder>
   ```
3. Lancez tous les services avec Docker Compose :
   ```bash
   docker-compose up --build
   ```
4. Accédez aux services :
   - **API** : [http://localhost:8080](http://localhost:8080)
   - **Web Client** : [http://localhost:3000](http://localhost:3000)
   - **Documentation** : [http://localhost:4000](http://localhost:4000)

### Documentation Contributeur
1. Clonez le projet :
   ```bash
   git clone <repository_url>
   cd <repository_folder>
   ```
2. Créez une branche pour vos changements :
   ```bash
   git checkout -b feature/<feature_name>
   ```
3. Respectez la convention de commit [Gitmoji](https://gitmoji.dev) pour documenter vos changements.
4. Testez localement :
   - Exécutez les tests unitaires :
     ```bash
     yarn test
     ```
   - Exécutez les tests E2E :
     ```bash
     yarn e2e
     ```
5. Ouvrez une Pull Request vers `develop` :
   - La PR doit inclure une description détaillée de vos changements.
   - Assurez-vous que les tests CI passent avant la fusion.
6. Après approbation, fusionnez votre PR.

---

## Rollback
Pour ce projet, j'ai choisi de gérer les rollbacks avec des **tags Git**. Cette méthode me permet :
1. D'identifier rapidement une version stable.
2. De déployer une version précédente en cas de problème.
3. D'assurer un suivi clair des versions.

---

## Résultats et Conclusion

Bien que je n'aie pas pu finaliser le déploiement sur GCP, ce projet a permis de mettre en place une infrastructure complète pour le développement, les tests, et le déploiement d'une application multi-services. Les étapes suivantes incluraient la résolution des problèmes liés à la connexion PostgreSQL et l'optimisation du déploiement Cloud Run.