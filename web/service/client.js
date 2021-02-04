import axios from 'axios';




const Client = () => {
  const options = {};  
  options.baseURL = process.env.NEXT_PUBLIC_API_URL;
  options.withCredentials = true;
  const instance = axios.create(options);  
  return instance;
};

export default Client;
