import React, { useState, useEffect } from 'react';
import Receipt, {isWarning} from '../models/receipt';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEllipsisV, faEye, faEdit, faTrash, faExclamationTriangle, faAngleLeft, faAngleRight, faAngleDoubleLeft, faAngleDoubleRight } from '@fortawesome/free-solid-svg-icons';
import { take, drop } from 'ramda';

interface ReceiptTableProps {
    receipts: Receipt[]
    itemsPerPage?: number
}

export default function ReceiptTable(props: ReceiptTableProps): JSX.Element {
    const { receipts, itemsPerPage } = props;
    const [page, setPage] = useState(0);

    const perPage = itemsPerPage ? itemsPerPage : 10;
    const pages = Math.ceil(receipts.length / perPage);

    return (
        <>
        <table className="xs:table-responsive-sm table table-sm text-sm lg:text-lg text-white">
            <thead>
                <tr className="bg-secondary">
                <th scope="col" className="fit"></th>
                <th scope="col">Vendor</th>
                <th scope="col">Date</th>
                <th scope="col">Total</th>
                <th scope="col">VAT</th>
                </tr>
            </thead>
            <tbody>
                {take(perPage, drop(page*perPage, receipts)).map((props) => {
                    return (
                        <Row key={props.id} needsReview={!props.verified} isWarning={isWarning(props)} {...props} />
                    );
                })}
            </tbody>
        </table>
        <div className="text-accent-2 italic">
            Showing {(page*perPage + 1).toString()} to {Math.min(receipts.length, ((page+1)*perPage)).toString()} of {receipts.length.toString()}
        </div>
        <div className="flex flex-row md:justify-start justify-between min-w-full text-accent-2">
                <div className="pl-0">
                    <button className="lg:px-2 lg:py-2 px-4 py-4" onClick={() => setPage(0)}><span><FontAwesomeIcon icon={faAngleDoubleLeft}/></span></button>
                </div>
                <div className="">
                    <button className="lg:px-2 lg:py-2 px-4 py-4" disabled={(page-1) < 0} onClick={() => setPage(page-1)}><span><FontAwesomeIcon icon={faAngleLeft}/></span></button>
                </div>
                <div className="px-3 lg:py-2 py-4">
                    Page {page + 1}
                </div>
                <div className="">
                    <button className="lg:px-2 lg:py-2 px-4 py-4" disabled={(page+1) == pages} onClick={() => setPage(page+1)}><span><FontAwesomeIcon icon={faAngleRight}/></span></button>
                </div>
                <div className="">
                    <button className="lg:px-2 lg:py-2 px-4 py-4" onClick={() => setPage(pages-1)}><span><FontAwesomeIcon icon={faAngleDoubleRight}/></span></button>
                </div>
        </div>
        </>
    );
}

interface RowProps extends Receipt {
    key: string
    isWarning: boolean
    needsReview: boolean
}

function Row(props: RowProps): JSX.Element {
    const [open, setOpen] = useState(false);
    
    // adds a window handler to detect clicks outside the dropdown to close
    useEffect(() => {
        const handleWindowClick = () => setOpen(false)
        if(open) {
          window.addEventListener('click', handleWindowClick);
        } else {
          window.removeEventListener('click', handleWindowClick)
        }
        return () => window.removeEventListener('click', handleWindowClick);
      }, [open, setOpen]);

    const { vendor_name, date, total, vat, isWarning, needsReview } = props;
    return (
        <tr scope="row" onClick={() => setOpen(!open)}>
            <td className="fit pr-5" onClick={() => setOpen(!open)}>
                <div className="pr-2 inline-block">
                    <div className="relative">
                        <button onClick={() => setOpen(!open)} className="relative z-10 block"><FontAwesomeIcon icon={faEllipsisV}/></button>
                        {open ? 
                            <div className="absolute left-0 mt-0 py-1 w-48 z-20 bg-secondary rounded-sm shadow-xl">
                                <div className="hover:bg-background">
                                    <a href="#" className="block px-4 py-2 text-white"><span className="pr-2 text-white"><FontAwesomeIcon icon={faEdit}/></span>View/Edit Receipt</a>
                                </div>
                                <div>
                                    <a href="#" className="block px-4 py-2 text-white"><span className="pr-3 text-red-700"><FontAwesomeIcon icon={faTrash}/></span>Delete Receipt</a>
                                </div>
                            </div> 
                        : null}
                    </div>           
                </div>
                {needsReview ? <span className="pl-1 text-accent-2"><FontAwesomeIcon icon={faEye}/></span> : null}
                {isWarning ? <span className="pl-1 text-accent-1"><FontAwesomeIcon icon={faExclamationTriangle}/></span> : null}
            </td>
            <td>{vendor_name}</td>
            <td>{date}</td>
            <td>{(total/100).toFixed(2)}</td>
            <td>{(vat/100).toFixed(2)}</td>
        </tr>
    );
}
