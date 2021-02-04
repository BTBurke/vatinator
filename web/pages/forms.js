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
import dayjs from 'dayjs';


export default function FormsPage() {
  const [account, setAccount] = useState({first_name: ''});
  const [editing, setEditing] = useState(false);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState(true);
  const [accountValid, setAccountValid] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const doAccount = async () => {
      await client().get('/api/account')
      .then((response) => {
        console.log(response);
        if (response.status === 200) {
          // returns 204 when there is no form data set,
          // only update if 200 with actual data
          setAccountValid(R.not(R.any(R.isEmpty)(R.values(response.data))));
          setAccount(response.data);
        }
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
  const updateAccount = (data) => {
    setAccountValid(R.not(R.any(R.isEmpty)(R.values(data))));
    const update = async (data) => await client().post('/api/account', data)
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
    update(data);
  };

  useEffect(() => setAccountValid(R.not(R.any(R.isEmpty)(R.values(account)))), [account]);
  
  //const accountValid = R.not(R.any(R.isEmpty)(R.values(account)));

  return (
    <>
    {loading ? 
       <Loading />
    :
    error ? 
      <Error error={error}/> 
    :
      <div className="container mx-auto">   
        <div className="w-full lg:w-3/4 mx-auto">
          <Nav />
          { message ? 
          <div className="py-0 bg-primary px-4 width-full rounded-md border-solid border-2 border-accent-2">
            <div className="ml-auto">
              <button className="p-2 text-white font-bold" onClick={() => setMessage(false)}>X</button>
            </div>
            <div className="px-2 py-1 w-full text-center text-white text-xl pb-4">
              For best results, be sure to read the <a className="underline" target="_blank" href="https://btburke.github.io/vatinator">documentation</a> first!
            </div>
          </div>
          : null }
          <div className="py-0 bg-primary px-4">
            {accountValid && !editing ? <Header>Create your forms</Header> : <Header>Enter your VAT form info</Header> }
            {accountValid && !editing ? <FormDetails account={account} setEditing={setEditing} showEdit /> : <AskForm initial={account} setEditing={setEditing} doUpdate={updateAccount} setAccount={setAccount} />}        
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

function Loading() {
  return (

        <p className="text-4xl text-gray-500 text-center w-full py-16">Loading...</p>
   
  );

}

function FileDrop(props) {
  const { onError } = props;
  const [doing, setDoing] = useState(null);
  const [rcpts, setRcpts] = useState(null);
  const [pct, setPct] = useState(0);
  const [batchID, setBatchID] = useState(null);
  const router = useRouter();
  
  useEffect(() => {
    setBatchID(Math.random().toString(16).substr(2, 14));
  }, []);

  const getBatchID = () => { 
    if (!batchID) {
      const batch = Math.random().toString(16).substr(2, 14);
      setBatchID(batch);
      return batch;
    } else {
      return batchID;
    }
  }


  const onDrop = (acceptedFiles) => {
    acceptedFiles.forEach((file) => {
      setDoing([`Uploading ${file.name}...`]);
      const reader = new FileReader()

      reader.onabort = () => onError('File reading was aborted')
      reader.onerror = () => onError('File reading has failed')
      reader.onload = async () => {
        let formdata = new FormData();
        formdata.append('file', file);
        formdata.append('name', file.name);
        await client().post('/api/file', formdata, 
          {
            params: {'batch_id': getBatchID()}, 
            headers: {'Content-Type': 'multipart/form-data'},
            onUploadProgress: event => {
                setPct(Math.round(100*event.loaded / event.total));
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
  };

  const {getRootProps, getInputProps, open, acceptedFiles} = useDropzone({
    // Disable click and keydown behavior
    noClick: true,
    noKeyboard: true,
    accept: ['image/*', 'application/zip', 'application/pdf'],
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
      <Process show={rcpts && rcpts !== 0} rcpts={rcpts} batchID={batchID} onError={onError} />
    </>
  );

}

function Process(props) {
  const submissionMonths = [
    dayjs().subtract(1, 'month').format('MMMM YYYY'),
    dayjs().format('MMMM YYYY')
  ]
  const defaultMonth = dayjs().date() <= 15 ? submissionMonths[0]: submissionMonths[1];

  const {show, rcpts, batchID, onError} = props;
  const router = useRouter();
  const [submissionDate, setSubmissionDate] = useState(defaultMonth);
  const [loading, setLoading] = useState(false);

  const handleProcess = async (e) => {
    e.preventDefault();
    console.log('submitting for processing ', batchID);
    setLoading(true);
    await client().post('/api/process', {batch_id: batchID, date: submissionDate})
    .then((response) => {
      if (response.status === 200) {
        router.push('/success');
      } else {
        onError(response.data)
      }
    })
    .catch((err) => {
      // don't turn off loading here to force page reload
      if (err.response) {
        onError(err.response.data);
      } else {
        onError('Something went wrong.')
      }
    });
  }

  if (!show) {
    return null;
  }
  return (
    <div className="w-full mx-auto py-0 pb-6 divide-y-2 divide-secondary">
      <div className="text-secondary font-bold py-1">
        Submit images for processing
      </div>
      <div className="py-4">
        <p className="text-lg text-gray-300">Select submission month</p>
        <select value={submissionDate} onChange={(e) => setSubmissionDate(e.target.value)} className="mt-1 py-1 appearance-none rounded bg-secondary text-white text-lg px-6 leading-tight">
          {submissionMonths.map((option) => (
                <option key={option} value={option}>{option}</option>
          ))}
        </select>
        <button disabled={loading} onClick={handleProcess} className="mt-6 bg-accent-2 w-full text-white px-full py-2 rounded-md font-bold border border-accent-2">
          {loading ? <span>Submitting...</span> :
          <>
            <span className="px-2"><FontAwesomeIcon icon={faCogs} /></span>  
            <span className="px-2">{rcpts === -1 ? `Process receipts` : `Process ${rcpts} receipts`}</span>
          </>
          }
          </button>
      </div>
    </div>
  )

}
