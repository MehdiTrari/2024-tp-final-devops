# Étape 1 : Build de la documentation
FROM node:18 as builder
WORKDIR /app

# Copier les fichiers de configuration
COPY package.json yarn.lock ./

# Installer les dépendances nécessaires
RUN yarn install

# Copier le reste du code
COPY . ./

# Construire la documentation
RUN yarn build

# Étape 2 : Image pour servir la documentation
FROM nginx:alpine
WORKDIR /usr/share/nginx/html

# Copier la documentation générée dans l'image finale
COPY --from=builder /app/build .

# Exposer le port 80
EXPOSE 80

# Commande de démarrage
CMD ["nginx", "-g", "daemon off;"]
