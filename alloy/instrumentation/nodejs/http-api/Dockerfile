FROM node:lts-slim
RUN mkdir -p /app/node_modules && chown -R node:node /app
WORKDIR /app
COPY package*.json ./
RUN chown node:node package*.json
USER node
RUN npm install
COPY --chown=node:node . .
EXPOSE 8080
CMD ["node", "index.js"]
