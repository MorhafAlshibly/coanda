FROM node:18

# Work directory
WORKDIR /

# Package JSON
COPY package*.json ./

# Install yarn
#RUN npm install yarn -g

# Install modules
RUN npm install

# Copy source files
COPY . .

# Build
RUN npm run build

# Expose to the API port
EXPOSE 5050

# Run
CMD ["node", "build/src/microservices/general/index.js"]