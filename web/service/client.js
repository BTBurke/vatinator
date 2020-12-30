import axios from 'axios';

const Client = () => {
  const options = {};  
  options.baseURL = 'http://localhost:8080';
  options.withCredentials = true;
  const instance = axios.create(options);  
  return instance;
};

export default Client;