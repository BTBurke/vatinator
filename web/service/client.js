import axios from 'axios';

const Client = () => {
  const options = {};  
  options.baseURL = 'https://api.vatinator.com';
  options.withCredentials = true;
  const instance = axios.create(options);  
  return instance;
};

export default Client;