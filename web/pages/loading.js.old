import React, { useState } from 'react';
import Loading from '../components/loading';

export default function LoadingPage() {
    const [isLoading, setLoading] = useState(false);

    return (
        <>
        {isLoading && <Loading />}
        <div className="container bg-base">
            <h1 className="text-accent-1">This is some content</h1>
            <div>
            <button className="text-white" onClick={() => setLoading(!isLoading)}>Set loading {!isLoading}</button>
            </div>
        </div>
        </>
        
    );

}