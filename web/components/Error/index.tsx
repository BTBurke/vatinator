import React from 'react';
import { faPooStorm } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';


export interface ErrorProps {
    error?: string
}

export default function Error(props: ErrorProps): JSX.Element {
    const {error} = props;
    return (
        <div className="top-0 left-0 bg-red-700 w-full absolute">
        <div className="flex flex-row content-center space-around justify-items-center items-center justify-center mx-auto max-w-xl">
            <div className="text-gray-100 px-4 py-4 text-5xl">
                <FontAwesomeIcon icon={faPooStorm}></FontAwesomeIcon>
            </div>
            <div className="px-4 py-4 text-gray-100 text-lg">
                <span className="font-bold">Error: </span>
                <span>{ error }</span>
                <div className="py-1">
                    Please reload the page and try again.
                </div>
            </div>
        </div>
        </div>
    );
}