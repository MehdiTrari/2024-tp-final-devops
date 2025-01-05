# Rapport de Déploiement avec Docker, Docker Compose et Intégration Continue

## Introduction
Dans ce projet, j'ai configuré et déployé une application multi-services composée de :
- Une API (`vote-api`),
- Une interface utilisateur (`web-client`),
- Une documentation statique (`docs`),
- Une base de données PostgreSQL (`db`).

J'ai utilisé Docker pour conteneuriser chaque service, Docker Compose pour orchestrer leur déploiement, et une pipeline d'intégration continue (CI) pour automatiser les tests et la validation du code.

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

Cette configuration garantit une validation continue du code et une gestion efficace des branches et des versions. La robustesse de la pipeline CI, la gestion agile des branches, et l'automatisation des déploiements avec Docker offrent une base solide pour le développement collaboratif et le déploiement fiable.
