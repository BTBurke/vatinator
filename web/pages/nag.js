import React  from 'react';
import Nav from '../components/nav';

export default function NagPage() {
    const barstyle = {
        width: '16%'
    }

    return (
        <>
        <div className="container mx-auto">
            <Nav/>
            <div className="w-full lg:w-3/4 mx-auto text-white bg-primary">
                <h1 className="text-2xl text-accent-1 lg:text-4xl font-bold py-2">Support the Vatinator</h1>
                <div className="py-1">
                    Running the Vatinator isn't free - <strong>it costs about $84/year</strong> for the server that makes this VAT form filling magic possible.

                    Now that your humble software craftsman has left Tallinn, will you help support the cost of keeping the Vatinator alive?
                </div>
                <div className="my-4 px-2 py-2 rounded-md bg-white">
                    <div className="flex justify-between mb-1">
                        <span className="text-sm font-medium text-black">$0</span>
                        <span className="text-sm font-medium text-black">$90†</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-3">
                        <div className="bg-blue-600 h-3 rounded-full" style={barstyle}></div>
                    </div>
                    <div className="flex justify-between mb-1">
                     <span></span>
                     <span className="text-xs font-medium text-black">† server costs + one beer</span>
                    </div>
                </div>
                <div className="inline-flex ml-auto">
                    <a href="/forms" type="button" className="mr-4 px-4 py-2 text-sm font-medium text-white bg-primary border border-white rounded-md">
                        Skip for this month
                    </a>
                    <a href="/forms" type="button" className="px-4 py-2 text-sm font-medium text-white bg-accent-2 border border-accent-2 rounded-md">
                        I already paid!
                    </a>
                </div>
                <div className="py-4">
                    <h1 className="text-lg text-accent-1 lg:text-xl font-bold py-2">How to Contribute</h1>
                    <div className="py-2">
                        Send some money to @BTBurke98 on Venmo... perhaps $10/year?
                    </div>
                </div>
            </div>
        </div>
        </>
    );
}
