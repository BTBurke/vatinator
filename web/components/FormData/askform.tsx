import React, {useState} from 'react';
import { FormBankData, FormBioData, FormEmbassyData, FormDetails } from './form';



export function AskForm(props): JSX.Element {
    const {setAccount, doUpdate, setEditing, initial} = props;
    const [account, setAccountLocal] = useState(initial || {
        first_name: "",
        last_name: "",
        diplomatic_id: "",
        embassy: "US Embassy",
        address: "Kentmanni 20",
        bank_name: "",
        account: "",
    });
    const [i, setI] = useState(0);

    const nextStep = data => {
        console.log(data);
        setAccountLocal(Object.assign(account, data))
        setI(i+1);
        console.log(account);
    }

    const submit = () => {
        setAccount(account);
        setEditing(false);
        doUpdate(account);
        console.log('final', account);
    }
    const reset = () => {
        setI(0);
    }

    const steps = [
        {
            child: <FormBioData onSubmit={nextStep} account={account} />,
        },
        {
            child: <FormEmbassyData onSubmit={nextStep} account={account} />,
        },
        {
            child: <FormBankData onSubmit={nextStep} account={account} />,
        },
        {
            child: <ConfirmFormDetails account={account} submit={submit} reset={reset} />,
        }
    ]
    const current = steps[i];
    return (
        <>
            {current.child}
        </>
    );
}

function ConfirmFormDetails(props): JSX.Element {
    const {account, submit, reset} = props;

    return (
        <>
            <FormDetails account={account} />
            <div>
                <p className="text-gray-500 text-lg font-bold py-2">Are these details correct?</p>
                <div>
                    <button onClick={submit} className="w-full lg:w-1/3 px-8 bg-accent-2 text-white py-2 rounded-md font-bold border border-accent-2">
                        Looks good
                    </button>
                </div>
                    
                <div className="py-4"> 
                    <button onClick={reset} className="w-full lg:w-1/3 px-8 bg-primary text-white py-2 rounded-md font-bold border border-white">
                        Or, go back and edit
                    </button>
                </div>
            </div>
        </>
    );
}