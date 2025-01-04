# Étape 1 : Build de l'application
FROM node:18 as builder
WORKDIR /app

# Copier les fichiers de configuration
COPY package.json yarn.lock ./
RUN yarn install

# Copier le reste du code
COPY . ./
RUN yarn build

# Étape 2 : Image pour l'exécution
FROM node:18-slim
WORKDIR /app

# Copier les fichiers nécessaires pour l'exécution
COPY --from=builder /app/package.json /app/yarn.lock ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

RUN yarn install --production

# Exposer le port
EXPOSE 3000

# Commande de démarrage
CMD ["yarn", "start"]
