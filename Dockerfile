FROM node:18-alpine

WORKDIR /app

COPY package.json ./
RUN yarn install
COPY . .
RUN yarn build
EXPOSE 5050

# Run
CMD ["node", "-r", "dotenv/config", "build/src/microservices/general/index.js"]
