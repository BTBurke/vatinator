import React, { useState } from 'react';

export default function Login() {
    const [email, setEmail] = useState(null);
    const [pw, setPw] = useState(null);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);

    const handleSubmit = () => {
        setLoading(true);
        console.log(email.value);
        console.log(pw.value);
        setInterval(() => {
            setError('Sorry, that didn\'t work');
            setLoading(false);
        }, 7000);
    }

    return (
        <div className="container m-auto">
            <div className="flex flex-col justify-center content-center items-center h-screen w-full">
                <div className="mx-auto pb-16">
                        <img src="\vatinatorAsPath.svg" className="h-20"></img>
                    </div>
                <div className="xs:w-75 md:w-50 xs:max-w-75 md:max-w-50">
                    <div className="p-4 bg-secondary rounded-md"> 
                        <label className="block">
                            <span className="text-accent-1 font-bold">Email</span>
                            <input ref={(input) => setEmail(input)} className="form-input mt-1 block w-full" id="email" placeholder="pompeo@state.gov"></input>   
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