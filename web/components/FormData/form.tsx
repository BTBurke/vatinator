import React, {useState, useEffect} from 'react';
import { useForm } from 'react-hook-form';
import { faEdit } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as R from 'ramda';

export function FormBioData(props): JSX.Element {
  const {onSubmit, account} = props;
  const { register, handleSubmit, watch, errors } = useForm();

  const localSubmit = data => {
    onSubmit(Object.assign(data, {full_name: data.first_name + ' ' + data.last_name}))
  }

  return (
    <form onSubmit={handleSubmit(localSubmit)}>
    <h2 className="text-secondary text-bold pb-1 text-2xl">Your personal data</h2>

      <div className="flex flex-col md:flex-row">
        <div className="w-full md:w-1/2 md:pr-2 py-2">
            <p className="text-gray-500 text-bold text-lg">First Name</p>
            <input name="first_name" defaultValue={account.first_name || ""} ref={register({required: true})} className="mt-1 py-1 appearance-none rounded bg-secondary text-white text-lg w-full leading-tight" />
            {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
        </div>
        <div className="w-full md:w-1/2 md:pl-2 py-2">
            <p className="text-gray-500 text-bold text-lg">Last Name</p>
            <input name="last_name" defaultValue={account.last_name || ""} ref={register({ required: true })} className="bg-secondary mt-1 py-1 appearance-none rounded text-white text-lg w-full leading-tight" />
            {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
        </div>
      </div>

    <div className="py-4">
      <p className="text-gray-500 text-bold text-lg">Diplomatic ID Number</p>
      <div>
        <input name="diplomatic_id" defaultValue={account.diplomatic_id || ""} ref={register({ required: true })} className="bg-secondary mt-1 py-1 appearance-none rounded text-white text-lg w-full md:w-1/2 leading-tight" />
      </div>
      {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
    </div>

      <div className="py-6">
        <input type="submit" value="Next" className="px-8 bg-accent-2 text-white py-2 rounded-md font-bold border border-accent-2"/>
      </div>
    
    </form>
  );
}

export function FormEmbassyData(props): JSX.Element {
    const {onSubmit, account} = props;
    const { register, handleSubmit, watch, errors } = useForm();
     
    return (
      <form onSubmit={handleSubmit(onSubmit)}>
      <h2 className="text-secondary text-bold pb-1 text-2xl">Your embassy data</h2>
  
        <div className="flex flex-col md:flex-row">
          <div className="w-full md:w-1/2 md:pr-2 py-2">
              <p className="text-gray-500 text-bold text-lg">Embassy Name</p>
              <input defaultValue={account.embassy || "US Embassy"} name="embassy" ref={register({required: true})} className="mt-1 py-1 appearance-none rounded bg-secondary text-white text-lg w-full leading-tight" />
              {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
          </div>
          <div className="w-full md:w-1/2 md:pl-2 py-2">
              <p className="text-gray-500 text-bold text-lg">Embassy Address</p>
              <input defaultValue={account.address || "Kentmanni 20"} name="address" ref={register({ required: true })} className="bg-secondary mt-1 py-1 appearance-none rounded text-white text-lg w-full leading-tight" />
              {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
          </div>
        </div>
  
        <div className="py-6">
          <input type="submit" value="Next" className="px-8 bg-accent-2 text-white py-2 rounded-md font-bold border border-accent-2"/>
        </div>
      
      </form>
    );
  }

  export function FormBankData(props): JSX.Element {
    const {onSubmit, account} = props;
    const { register, handleSubmit, watch, errors } = useForm();
    const localSubmit = data => {
      console.log(data);
        const bank_name_computed = data.bank_name !== 'other' ? data.bank_name : data.bank_name_other;
        onSubmit(Object.assign(data, {bank: bank_name_computed + ', ' + data.account}))
    }
    const [showOther, setShowOther] = useState(false);

    useEffect(() => {
        if (watch('bank_name') === 'other') {
        setShowOther(true);
      } else {
        setShowOther(false);
      }
    }, [watch('bank_name')]); 
      
    return (
      <form onSubmit={handleSubmit(localSubmit)}>
      <h2 className="text-secondary text-bold pb-1 text-2xl">Your bank data</h2>
  
        <div className="flex flex-col md:flex-row">
          <div className="w-full md:w-1/2 md:pr-2 py-2">
              <p className="text-gray-500 text-bold text-lg">Bank</p>
              <select defaultValue={account.bank_name || ""} name="bank_name" ref={register({required: true})} className="mt-1 py-1 appearance-none rounded bg-secondary text-white text-lg w-full leading-tight">
                  <option value="">-- Select your bank --</option>
                  <option value="AS SEB Bank, EEUHEE2X, Tornimäe 2, 15010 Tallinn, Estonia">SEB Bank</option>
                  <option value="Swedbank AS, HABAEE2X, Liivalaia 8, 15040 Tallinn, Estonia">Swedbank</option>
                  <option value="AS LHV Bank, LHVBEE22, Tartu mnt 2, 10145 Tallinn, Estonia">LHV Bank</option>
                  <option value="Luminor Bank AS, NDEAEE2X, Liivalaia 45, 10145 Tallinn, Estonia">Luminor Bank</option>
                  <option value="other">Some other bank</option>
              </select>
              {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
          </div>

          <div className="w-full md:w-1/2 md:pl-2 py-2">
              <p className="text-gray-500 text-bold text-lg">Account Number</p>
              <input defaultValue={account.account || ""} name="account" ref={register({ required: true })} className="bg-secondary mt-1 py-1 appearance-none rounded text-white text-lg w-full leading-tight" />
              {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
          </div>
        </div>
        {showOther ? <div className="w-full py-2">
              <p className="text-gray-500 text-bold text-lg">Bank SWIFT Code and Address</p>
              <p className="text-secondary text-md">(e.g., Swedbank AS, HABAEE2X, Liivalaia 8, 15040 Tallinn, Estonia)</p>
              <input defaultValue={account.bank_name || ""} name="bank_name_other" ref={register({ required: false })} className="bg-secondary mt-1 py-1 appearance-none rounded text-white text-lg w-full leading-tight" />
              {errors.first_name && <span className="text-red-800 text-bold">This field is required</span>}
          </div> : null}
  
        <div className="py-6">
          <input type="submit" value="Save" className="px-8 bg-accent-2 text-white py-2 rounded-md font-bold border border-accent-2"/>
        </div>
      
      </form>
    );
  }

export function FormDetails(props): JSX.Element {
    const { account, setEditing, showEdit } = props;
    const accountValid = R.not(R.any(R.isEmpty)(R.values(account)));
    
    if (!accountValid) {
        return (
            <>
            <p className="text-white">Loading...</p>
            </>
        ) 
    }
    return (
    <div className="divide-y-2 divide-secondary">
      <div className="flex flex-row justify-between mt-2">
          <div className="text-secondary font-bold py-1">
              Form Info
          </div>
          {showEdit && <div className="text-secondary text-md px-3 py-1" onClick={() => setEditing(true)}>
              <FontAwesomeIcon icon={faEdit}></FontAwesomeIcon>
          </div>}
      </div>
      <div className="pt-2 flex md:flex-row flex-col justify-between mx-auto lg:px-0 xs:px-0 pb-2 md:w-full">
              <div className="text-white text-left">
                <p className="text-sm font-bold lg:text-lg">{account.full_name} ({account.diplomatic_id})</p>
                <p className="text-sm lg:text-lg">
                  <span>{account.embassy}, {account.address}</span>
                </p>
              </div>
              <div className="text-white py-4 md:py-0">
                <p className="text-sm font-bold lg:text-lg">{account.bank.split(',').slice(0,2).join(',')}</p>
                <p className="text-sm lg:text-lg">{account.bank.split(',').slice(2, account.bank.length-1).join(',')}</p>
              </div>
      </div>
    </div> 
    );
}