import React, { useState, useEffect } from 'react';
import client from '../service/client';
import { useRouter } from 'next/router'; 

export default function Login() {
    const [email, setEmail] = useState(null);
    const [pw, setPw] = useState(null);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const router = useRouter();
    const [checking, setChecking] = useState(true);

    useEffect(() => {
        client().get('/session')
        .then((response) => {
            setChecking(false);
            if (response.status === 200) {
                router.push('/nag');
            }
        })
        .catch((err) => { 
            setChecking(false);
            console.log('no existing session'); 
        });
    }, []);
    useEffect(() => {
        setTimeout(() => setChecking(false), 450);
    }, []);

    const handleSubmit = () => {
        setLoading(true);
        if (pw.value.length === 0 || email.value.length === 0) {
            setError('Email and password are required.');
            setLoading(false);
            return
        }

        const doLogin = async (email, password) => await client().post('/login',{email: email, password: password})
        .then((response) => {
            if (response.status === 200) {
                setLoading(false);
                router.push('/forms');
            }
        })
        .catch((err) => {
            setLoading(false);
            if (err.response) {
                setError(err.response.data);
            } else {
                setError('Sorry, that didn\'t work');
            }
        });
        doLogin(email.value, password.value);
    }
    if (checking) {
        return null;
    }
    return (
        <div className="container m-auto">
            <div className="flex flex-col justify-center content-center items-center h-screen w-full">
                <div className="mx-auto pb-16">
                        <img src="\vatinatorAsPath.svg" className="h-20"></img>
                    </div>
                <div className="xs:w-75 md:w-50 xs:max-w-75 md:max-w-50 pb-16">
                    <h1 className="text-2xl font-bold text-white">Login</h1>
                    <div className="p-4 bg-secondary rounded-md"> 
                        <label className="block">
                            <span className="text-accent-1 font-bold">Email</span>
                            <input ref={(input) => setEmail(input)} className="form-input mt-1 block w-full" id="email" placeholder="email"></input>   
                        </label>
                        <label className="block pt-5">
                            <span className="text-accent-1 font-bold">Password</span>
                            <input ref={(input) => setPw(input)} className="form-input mt-1 block w-full" id="password" type="password" placeholder="password"></input>   
                        </label>
                        <div className="pt-1">
                            <a className="text-white italic" href="mailto:burkebt@state.gov?subject=reset%20password">Forgot password?</a>
                        </div>
                        <div className="pt-6">
                            <button disabled={loading} onClick={handleSubmit} className="w-full bg-accent-2 p-2 rounded-full text-white font-bold">{loading ? "Logging in..." : "Login"}</button>
                        </div>
                        <div className="text-center break-words text-red-700 italic font-bold">
                            {error ? error : null}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

}