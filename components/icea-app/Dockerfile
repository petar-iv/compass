from node:16
# Create app directory
WORKDIR /usr/node/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD [ "node", "index.js" ]