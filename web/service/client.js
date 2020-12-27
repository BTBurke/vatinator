import axios from 'axios';

const Client = () => {
  const options = {};  
  options.baseURL = 'http://localhost:8080';
  const instance = axios.create(options);  
  return instance;
};

export default Client;