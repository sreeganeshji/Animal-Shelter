import React, { useEffect, useState } from "react";

import { getData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const ExpenseAnalysis = () => {

    const [state, setState] = useState({
        data: [],
        errorMessage: "",
        loginRedirect: false,
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer) {
            setState({ ...state, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    useEffect(() => {
        if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer) {
            getData("/expense-analysis")
                .then(({ target }) => {
                    if (target.status === 200) {
                        setState({ ...state, data: JSON.parse(target.response) });
                    } else {
                        setState({ ...state, errorMessage: target.response });
                    }
                })
                .catch(error => {
                    console.log(error);
                    setState({ ...state, errorMessage: JSON.stringify(error) });
                })
        }
    }, [state.cookie]);
    return (
        <>
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Expense Analysis</h1>
                {

                    state.data.length ?
                        // && (return <div className="form-unit">{JSON.stringify(datum)}</div>

                        <table>
                            <thead>
                                <tr>
                                    <th>Vendor</th>
                                    <th>Total Spending in dollars</th>
                                </tr>
                            </thead>
                            <tbody>
                                {
                                    state.data.map((data, index) => {
                                        return (<tr>
                                            <td>
                                                {data.vendor}
                                            </td>
                                            <td>
                                                ${data.totalSpending / 100}
                                            </td>
                                        </tr>)
                                    })
                                }
                            </tbody>
                        </table>

                        :
                        <p>No expenses found</p>

                }

                {
                    state.errorMessage && (
                        <div styles={{ color: "red" }}>{state.errorMessage}</div>
                    )}
            </main>
        </>
    );
}

export default ExpenseAnalysis;