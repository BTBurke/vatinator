import React from 'react';
import { NavHome } from '../components/nav';
import '@fontsource/bungee-shade';

export default function IndexPage() {

    return (
        
        <div className="container mx-auto">   
            <div className="w-full lg:w-3/4 mx-auto">
                <NavHome />
            </div>

            <h1 className="pt-10 font-display text-4xl md:text-6xl text-center md:text-left text-white w-full lg:w-3/4 mx-auto" style={{fontFamily: 'Bungee Shade'}}>
                You've found<br></br>the <span className="text-accent-1">Vatinator</span>
            </h1>
            <div className="block md:hidden mx-auto text-center pt-24">
                <a className="mr-4 bg-primary text-white px-full py-2 px-6 rounded-md font-bold border border-white" href="/login">Login</a>
                <a className="mr-4 bg-primary text-white px-full py-2 px-6 rounded-md font-bold border border-white" href="/create">Create account</a>
            </div>

            <div className="pt-10 md:pt-8 md:text-2xl md:text-left text-white text-lg text-center px-1 md:px-0 lg:w-3/4 mx-auto w-full">
              Have questions? Check out the <a className="underline" href="https://btburke.github.io/vatinator/">documentation</a> to get started.
            </div>
        </div>
    
    );
}
