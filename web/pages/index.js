import React, { useCallback, useState }  from 'react';
import Nav from '../components/nav'
import { faReceipt, faFolderOpen} from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useDropzone } from 'react-dropzone';
import Loading from '../components/loading';

export default function IndexPage() {
  const [doing, setDoing] = useState(null); 
  
  // TODO: add a application/zip file handler
  
  const onDrop = useCallback((acceptedFiles) => {
    if (acceptedFiles.length > 1) {
      // create a batchProcess
      console.log(`total files ${acceptedFiles.length}`);
    } 
    acceptedFiles.forEach((file) => {
      setDoing([`Processing ${file.name}...`]);
      console.log(file);
      const reader = new FileReader()

      reader.onabort = () => console.log('file reading was aborted')
      reader.onerror = () => console.log('file reading has failed')
      reader.onload = () => {
      // Do whatever you want with the file contents
        const binaryStr = reader.result;
        console.log(`got result ${binaryStr}`);
      }
      reader.readAsArrayBuffer(file);
    })

    
  }, []);

  const {getRootProps, getInputProps, open, acceptedFiles} = useDropzone({
    // Disable click and keydown behavior
    noClick: true,
    noKeyboard: true,
    accept: ['image/*'],
    onDrop,
  });

  return (
    <div>
      { doing && <Loading msgs={doing} /> }
      <div className="lg:container lg:mx-auto">
        <Nav />
        <div className="py-0 bg-primary px-4">
          <p className="text-2xl text-accent-1 lg:text-4xl font-bold">
                Current Batch
          </p>
          <div className="flex flex-row justify-between lg:px-0 xs:px-0 py-2">
            <div className="text-white text-left">
              <p className="text-sm lg:text-lg">STARTED</p>
              <p className="text-lg font-bold lg:text-xl">
                <span>Oct 21</span>
                <span className="text-sm px-2">(21 days ago)</span>
              </p>
            </div>
            <div className="text-white text-center">
              <p className="text-sm lg:text-lg">RECEIPTS</p>
              <p className="text-lg font-bold lg:text-xl">60</p>
            </div>
            <div className="text-white text-center">
              <p className="text-sm lg:text-lg">REFUND</p>
              <p className="text-lg font-bold lg:text-xl">160.00€</p>
            </div>
          </div>
          
          
          
            <div {...getRootProps()} className="md:w-full lg:w-3/4 mx-auto my-10 md:py-16 md:px-16 md:min-h-1/2 md:border-dashed md:border-secondary md:border-2 md:rounded-sm">
                   <input {...getInputProps()}></input>
                   <p className="hidden md:block text-secondary text-center italic pb-2 font-bold">You can drop image files here or click to select receipt(s)</p>
                  <button onClick={open} className="bg-accent-2 w-full text-white px-full py-2 md:mb-2 rounded-full font-bold border border-accent-2">
                    <span className="px-2"><FontAwesomeIcon icon={faReceipt} /></span>  
                    <span className="px-2">Add receipt</span>
                  </button>
              
            </div>
            <div className="md:w-full lg:w-3/4 mx-auto py-0">
              <button className="bg-primary w-full text-white px-full py-2 rounded-full font-bold border border-white">
                <span className="px-2"><FontAwesomeIcon icon={faFolderOpen} /></span>  
                <span className="px-2">Manage receipts</span>
              </button>
            </div>
  
        </div>
      </div>
    </div>
  )
}
