import React, { useCallback, useEffect, useState }  from 'react';
import Nav from '../components/nav'
import { faReceipt, faCogs} from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useDropzone } from 'react-dropzone';
import axios from 'axios';
import { AskForm , FormDetails } from '../components/FormData';
import Error from '../components/Error';
import * as R from 'ramda';
import client from '../service/client';
import { useRouter } from 'next/router'

export default function IndexPage() {
  const [account, setAccount] = useState({leave_this_here: ''});
  const [editing, setEditing] = useState(false);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const doAccount = async () => {
      await client().get('/account')
      .then((response) => {
        console.log(response);
        setAccount(response.data);
        setLoading(false);
      })
      .catch((err) => {
        setLoading(false);
        if (err.response) {
          if (err.response.status === 401 || err.response.status === 403) {
            router.push('/login');
          } else {
            setError(err.response.data);
          }
        } else {
          setError('Something went wrong.');
        }
      });
      return;
    };
    doAccount();
  }, []);
  const updateAccount = async (data) => {
    setEditing(false);
    await client().post('/account', data)
    .then((response) => {
      if (response.status === 200) {
        setAccount(data);
      } else {
        throw new Error("unknown response")
      }
    })
    .catch((err) => {
      if (err.response) {
        setError(err.response.data);
      } else {
        setError('Something went wrong updating account data.')
      }
    });
  };
  
  const accountValid = R.not(R.any(R.isEmpty)(R.values(account)));

  return (
    <>
    {loading ? 
      <p className="text-4xl text-gray-500 text-center w-full py-16">Loading...</p> 
    :
    error ? 
      <Error error={error}/> 
    :
      <div className="container mx-auto">   
        <div className="w-full lg:w-3/4 mx-auto">
          <Nav />
          <div className="py-0 bg-primary px-4">
            {accountValid && !editing ? <Header>Create your forms</Header> : <Header>Enter your VAT form info</Header> }
            {accountValid && !editing ? <FormDetails account={account} setEditing={setEditing} showEdit /> : <AskForm initial={account} setAccount={updateAccount} />}        
            {accountValid && !editing ? <FileDrop onError={setError} /> : null }   
          </div>
        </div>
      </div>
    }
    
    </>
  )
}

function Header(props) {
  return (
    <p className="text-2xl text-accent-1 lg:text-4xl font-bold">
        {props.children}
    </p>
  );
}

function FileDrop(props) {
  const { onError } = props;
  const [doing, setDoing] = useState(null);
  const [rcpts, setRcpts] = useState(null);
  const [pct, setPct] = useState(0);
  const [batchID, setBatchID] = useState("");
  const router = useRouter();
  
  useEffect(() => {
    setBatchID(Math.random().toString(16).substr(2, 14));
  }, []);


  const onDrop = useCallback((acceptedFiles) => {
    acceptedFiles.forEach((file) => {
      setDoing([`Processing ${file.name}...`]);
      const reader = new FileReader()

      reader.onabort = () => onError('File reading was aborted')
      reader.onerror = () => onError('File reading has failed')
      reader.onload = async () => {
        let formdata = new FormData();
        formdata.append('file', file);
        formdata.append('name', file.name);
        await client().post('/file', formdata, 
          {
            params: {'batch_id': batchID}, 
            headers: {'Content-Type': 'multipart/form-data'},
            onUploadProgress: event => {
              setPct(Math.round((100 * event.loaded) / event.total));
            },
        })
        .then(() => {
          console.log('uploaded ', file.name); 
          setDoing(null);
          if (file.name.endsWith('.zip')) {
            // set sentinel value for zip with unknown number of files, maybe
            // could return number of files here which would be good
            setRcpts(-1);
          } else {
            // if already set to sentinel -1 because of a zip file, just leave it.  Otherwise,
            // count number of receipts uploaded
            if (rcpts !== -1) {
              setRcpts(rcpts+acceptedFiles.length);
            }
          }
        }).catch((err) => {
          if (err.response) {
            onError(err.response.data);
          } else {
            onError('Something went wrong.');
          }
        });
        return
      }
      reader.readAsBinaryString(file);
    }) 
  }, []);

  const handleProcess = (e) => {
    e.preventDefault();
    console.log('submitting for processing ', batchID);
    router.push('/success');
  }

  const {getRootProps, getInputProps, open, acceptedFiles} = useDropzone({
    // Disable click and keydown behavior
    noClick: true,
    noKeyboard: true,
    accept: ['image/*', 'application/zip'],
    onDrop,
  });

  return (
    <>
      <div {...getRootProps()} className="md:w-full mx-auto my-10 md:py-16 md:px-16 md:min-h-1/2 md:border-dashed md:border-secondary md:border-2 md:rounded-sm">
        <input {...getInputProps()}></input>
        {!doing && <p className="hidden md:block text-secondary text-center italic pb-2 font-bold">You can drop images or a zip file here or click to select receipt(s)</p>}
        {doing ? <p className="block text-2xl text-gray-500 text-center italic pb-2 font-bold">{`Uploading...${pct}%`}</p> : 
        
        <button onClick={open} className={rcpts ? "bg-primary w-full text-white px-full py-2 md:mb-2 rounded-md font-bold border border-white" : "bg-accent-2 w-full text-white px-full py-2 md:mb-2 rounded-md font-bold border border-accent-2"}>
          <span className="px-2"><FontAwesomeIcon icon={faReceipt} /></span>  
          <span className="px-2">{rcpts ? `Add more receipts` : `Add receipts`}</span>
        </button> }
      </div>
        {rcpts && rcpts !== 0 ? 
        <div className="md:w-full lg:w-3/4 mx-auto py-0 pb-6">
          <button onClick={handleProcess} className="bg-accent-2 w-full text-white px-full py-2 rounded-md font-bold border border-accent-2">
            <span className="px-2"><FontAwesomeIcon icon={faCogs} /></span>  
            <span className="px-2">{rcpts === -1 ? `Process receipts` : `Process ${rcpts} receipts`}</span>
          </button>
        </div> : null}
    </>
  );

}