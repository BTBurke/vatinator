import React, { useEffect, useState } from 'react';
import styles from '../styles/loading.module.scss';

export default function Loading(props) {
    const { msgs } = props;
    const msg = msgs ? msgs : ['Working...', 'Still working...', 'Please be patient...'];
    
    const [m, setM] = useState(0);
    

    useEffect(() => {
        setTimeout(() => {
            if (m === (msg.length-1)) {
                return
            } else {
                setM(m+1);    
            }
        }, 3000);
    }, [m]);

    return (
        <div className={styles.ov}>
            <div>
                <div className="m-auto px-2 py-2 bg-secondary w-24 h-24 rounded-lg">
                    <div className={styles.lds}><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div></div>
                </div>
                <div className="text-white text-center w-full font-bold mt-4">
                    {msg[m]}
                </div>
            </div>
        </div>
    );

}