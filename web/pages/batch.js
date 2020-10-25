import React, { useState } from 'react';
import Nav from '../components/nav';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTasks, faFilePdf } from '@fortawesome/free-solid-svg-icons';
import ReceiptTable from '../components/table';
import { isWarning , isVerified } from '../models/receipt';
import { filter } from 'ramda';

export default function Batch() {
    const testData = [
        {verified: new Date(), id: "1", vendor_name: "Bauhof Group AS 1", date: "09/27", tax_id: "test", receipt_number: "EEtest", total: 1300,  vat: 500},
        {id: "2", vendor_name: "Bauhof Group AS 2", date: "09/27", total: 1300,  vat: 500},
        {id: "3", vendor_name: "Bauhof Group AS 3", date: "09/27", total: 1300,  vat: 500},
        {id: "4", vendor_name: "Bauhof Group AS 4", date: "09/27", total: 1300,  vat: 500},
        {id: "5", vendor_name: "Bauhof Group AS 5", date: "09/27", total: 1300,  vat: 500},
        {id: "6", vendor_name: "Bauhof Group AS 6", date: "09/27", total: 1300,  vat: 500},
        {id: "7", vendor_name: "Bauhof Group AS 7", date: "09/27", total: 1300,  vat: 500},
        {id: "8", vendor_name: "Bauhof Group AS 8", date: "09/27", total: 1300,  vat: 500},
        {id: "9", vendor_name: "Bauhof Group AS 9", date: "09/27", total: 1300,  vat: 500},
        {id: "10", vendor_name: "Bauhof Group AS 10", date: "09/27", total: 1300,  vat: 500},
        {id: "11", vendor_name: "Bauhof Group AS 11", date: "09/27", total: 1300,  vat: 500},
        {id: "12", vendor_name: "Bauhof Group AS 12", date: "09/27", total: 1300,  vat: 500},
        {id: "13", vendor_name: "Bauhof Group AS 13", date: "09/27", total: 1300,  vat: 500},
        {id: "14", vendor_name: "Bauhof Group AS 14", date: "09/27", total: 1300,  vat: 500},
        {id: "15", vendor_name: "Bauhof Group AS 15", date: "09/27", total: 1300,  vat: 500},
        {id: "16", vendor_name: "Bauhof Group AS 16", date: "09/27", total: 1300,  vat: 500},
        {id: "17", vendor_name: "Bauhof Group AS 17", date: "09/27", total: 1300,  vat: 500},
        {id: "18", vendor_name: "Bauhof Group AS 18", date: "09/27", total: 1300,  vat: 500},
        {id: "19", vendor_name: "Bauhof Group AS 19", date: "09/27", total: 1300,  vat: 500},
        {id: "20", vendor_name: "Bauhof Group AS 20", date: "09/27", total: 1300,  vat: 500},
        {id: "21", vendor_name: "Bauhof Group AS 21", date: "09/27", total: 1300,  vat: 500},
        {id: "22", vendor_name: "Bauhof Group AS 22", date: "09/27", total: 1300,  vat: 500},
    ]
    // const testData = [
    //     {verified: new Date(), id: "1", vendor_name: "Bauhof Group AS 1", date: "09/27", tax_id: "test", receipt_number: "EEtest", total: 1300,  vat: 500},
    // ]

    const forReview = filter((r) => { return isWarning(r) || !isVerified(r) }, testData);

    return (
      <div className="lg:container lg:mx-auto">
        <Nav />
        <div className="py-0 bg-primary px-4">
            <p className="text-2xl text-accent-1 lg:text-4xl font-bold">
                Batch
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
        </div>
        <div className="px-4">
            <p className="text-2xl text-accent-1 lg:text-4xl font-bold pt-4">
                Receipts
            </p>
            <div className="pt-2">
                <ReceiptTable receipts={testData} />
            </div>
        </div>
        <div className="md:w-full lg:w-3/4 mx-auto py-4 px-2">
                    { forReview.length > 0 ?
                        <button className="bg-accent-2 w-full text-white px-full py-2 rounded-full font-bold border border-accent-2">
                            <span className="px-2"><FontAwesomeIcon icon={faTasks} /></span>  
                            <span className="px-2">Review {forReview.length.toString()} receipts with issues</span>
                        </button>
                    :
                        <button className="bg-accent-2 w-full text-white px-full py-2 rounded-full font-bold border border-accent-2">
                            <span className="px-2"><FontAwesomeIcon icon={faFilePdf} /></span>  
                            <span className="px-2">Close batch and create VAT submission</span>
                        </button>
                    }
        </div>
      </div>
    );
}
