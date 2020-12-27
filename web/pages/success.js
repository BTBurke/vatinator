import React from 'react';
import { faThumbsUp } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export default function SuccessPage() {
    return (
        <div className="py-4 max-w-lg container m-auto">
            <div className="flex flex-col lg:flex-row bg-primary px-6 py-4 justify-center content-center items-center h-screen w-full">
                <div className="text-6xl text-accent-2">
                    <FontAwesomeIcon icon={faThumbsUp}></FontAwesomeIcon>
                </div>
                <div className="text-xl px-4">
                    <span className="text-accent-1 font-bold mr-2">Success!</span>
                    <span className="text-white text-lg">I'm working on your forms.  You'll receive an email with a download link when they are ready.  See you next month.</span>
                </div>
            </div>
        </div>
    );
}