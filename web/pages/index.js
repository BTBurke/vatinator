import React from 'react';
import { NavHome } from '../components/nav';

export default function IndexPage() {

    return (
        
        <div className="container mx-auto">   
            <div className="w-full lg:w-3/4 mx-auto">
                <NavHome />
            </div>

            <h1 className="pt-10 font-display text-6xl text-white w-full lg:w-3/4 mx-auto">
                You've found<br></br>the <span className="text-accent-1">Vatinator</span>
            </h1>
        </div>
    
    );
}