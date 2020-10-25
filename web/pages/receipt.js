import React from 'react';
import { useState, useEffect } from 'react';

import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";
import Nav from '../components/nav';
import styles from '../styles/receipt.module.scss';
import { faThumbsUp, faTrash, faSave } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useRouter } from 'next/router';


export default function ReceiptPage() {
    // data structure for edits, if property exists here, it has been clicked to edit
    const [editing, setEditing] = useState(null);

    // router to get intent from query
    const router = useRouter();
    const intent = router.query.intent;
    
    // set up hierarcy of intent, one of delete, review, or other
    // review will save changes and progress to the next unreviewed receipt
    // expects query like /receipt/[id0]?intent=review&review=id1,id2,id3
    const handleIntent = (localIntent) => {
        return function() {
            if (localIntent === 'delete') {
                // do delete stuff
                console.log('would delete');
            } else if (intent === 'review') {
                // do review stuff, save current then push to next in list
                console.log('would go to next in list');
            } else {
                // do save then return to home
                console.log('would save current, look for diff');
            }
        }
    };
    
    // controls focus when multiple items are edited
    const [ref, setRef] = useState(null);
    useEffect(() => {
        if (ref !== null) {
            ref.focus();
        }
        setRef(null);
    });

    const rcpt = {
        vendor_name: 'Bauhof Group AS',
        date: '09/27/2020',
        receipt_number: '1010500205581',
        tax_id: 'EE100589638',
        total: '94.74',
        vat: '15.79',
    }
    // TODO: it should get a dynamically-generated responsive image size that fits in the viewport
    // like https://images.vatinator.com/receipt/id?height=714
    // then get rid of all of the overflow handling
    return (
        <div className="">
            <div className="flex flex-col lg:container lg:mx-auto max-h-screen">
                <div className="flex-shrink px-8 py-4 mx-auto overflow-y-scroll overflow-x-scroll object-contain max-h-2/3vh">
                        <TransformWrapper options={{limitToWrapper: true}} reset={{disabled: true}}>
                            <TransformComponent>
                                <img className="h-full object-fill" src="/test2.jpg" alt="test" />
                            </TransformComponent>
                        </TransformWrapper>
                </div>
                <div className="max-h-1/3vh px-8 py-4 h-1/4vh"> 
                    <div className={styles.container}>
                            <Field field="vendor" label="VENDOR" value={rcpt.vendor_name} editing={editing} setEditing={setEditing} setRef={setRef} />
                            <Field field="date" label="DATE" value={rcpt.date} editing={editing} setEditing={setEditing} setRef={setRef}/>
                            <Field field="id" label="RECEIPT NUMBER" value={rcpt.receipt_number} editing={editing} setEditing={setEditing} setRef={setRef} />
                            <Field field="tax_id" label="TAX ID" value={rcpt.tax_id} editing={editing} setEditing={setEditing} setRef={setRef} />
                            <Field field="total" label="TOTAL" value={rcpt.total} editing={editing} setEditing={setEditing} setRef={setRef} />
                            <Field field="vat" label="VAT" value={rcpt.vat} editing={editing} setEditing={setEditing} setRef={setRef} />

                            <div className={styles.btn}>
                                {editing ?
                                    <button onClick={handleIntent('save')} className="bg-accent-2 w-full text-white px-full py-2 rounded-full font-bold border border-accent-2">
                                        <span className="px-2"><FontAwesomeIcon icon={faSave} /></span>  
                                        <span className="px-2">Save changes</span>
                                    </button>
                                : intent === 'delete' ?
                                    <button onClick={handleIntent('delete')} className="bg-red-700 w-full text-white px-full py-2 rounded-full font-bold border border-red-700">
                                        <span className="px-2"><FontAwesomeIcon icon={faTrash} /></span>  
                                        <span className="px-2">Permanently delete receipt</span>
                                    </button>
                                :
                                    <button onClick={handleIntent('save')} className="bg-accent-2 w-full text-white px-full py-2 rounded-full font-bold border border-accent-2">
                                        <span className="px-2"><FontAwesomeIcon icon={faThumbsUp} /></span>  
                                        <span className="px-2">Everything looks good</span>
                                    </button>    
                                }                               
                                <div className="text-xs text-center text-white">or click on any value to edit</div>
                            </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

function Field(props) {

    var localRef = null;

    const { field, label, value, editing, setEditing, setRef } = props;
    
    // adds property to the editing data then controls focus through ref
    function handleDivClick() {
        if (!editing || editing[field] === undefined) {
         setEditing(Object.assign({}, editing, {[`${field}`]: value}));
        }
        setRef(localRef);
    }


    // input updates state on change, final values of edited items are in state
    return (
        <div className={styles[field]} onClick={() => handleDivClick()}>
            <div className="text-xs text-accent-1 font-bold">{label}</div>
            <div className="text-white">
                {editing && editing[field] !== undefined ? <input ref={(input) => { localRef=input; setRef(input);} } className="bg-secondary appearance-none rounded text-white w-full leading-tight" id={field} type="text" value={editing[field]} onChange={() => setEditing(Object.assign({}, editing, {[`${field}`]: localRef.value}))}></input> : value }
            </div>
        </div>
    );


}