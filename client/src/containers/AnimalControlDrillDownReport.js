import React, { useEffect, useState } from "react";
import { Redirect, useParams } from "react-router-dom";

import { getData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const AnimalControlDrillDownReport = () => {

    const [state, setState] = useState({
        errorMessage: "",
        cookie: {},
        backRedirect: false
    });

    const [surrenderedData, setSurrenderedData] = useState([]);

    const [adoptedData, setAdoptedData] = useState([]);

    const { monthAndYear } = useParams();

    const [month, year] = monthAndYear.split("-");

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
            const startMonth = (() => {
                switch (month) {
                    case "January":
                        return "01";
                    case "February":
                        return "02";
                    case "March":
                        return "03";
                    case "April":
                        return "04";
                    case "May":
                        return "05";
                    case "June":
                        return "06";
                    case "July":
                        return "07";
                    case "August":
                        return "08";
                    case "September":
                        return "09";
                    case "October":
                        return "10";
                    case "November":
                        return "11";
                    case "December":
                        return "12";
                    default:
                        throw new Error("Invalid month provided");
                }
            })();
            const endMonth = (() => {
                switch (month) {
                    case "January":
                        return "02";
                    case "February":
                        return "03";
                    case "March":
                        return "04";
                    case "April":
                        return "05";
                    case "May":
                        return "06";
                    case "June":
                        return "07";
                    case "July":
                        return "08";
                    case "August":
                        return "09";
                    case "September":
                        return "10";
                    case "October":
                        return "11";
                    case "November":
                        return "12";
                    case "December":
                        return "01";
                    default:
                        throw new Error("Invalid month provided");
                }
            })();

            const endYear = (() => {
                return month === "December" ? year + 1 : year;
            })();

            getData(`/animal-control-report-drilldown-surrendered?startDate=${year + "-" + startMonth + "-01"}&endDate=${endYear + "-" + endMonth + "-01"}`)
                .then(({ target }) => {
                    const surrenderedData = JSON.parse(target.response);
                    if (target.status === 200 && surrenderedData.length) {
                        setSurrenderedData(surrenderedData);
                    }
                })
                .catch(error => {
                    console.log(error);
                    setState({ ...state, errorMessage: error });
                });

            getData(`/animal-control-report-drilldown-adopted?startDate=${year + "-" + startMonth + "-01"}&endDate=${endYear + "-" + endMonth + "-01"}`)
                .then(({ target }) => {
                    const adoptedData = JSON.parse(target.response);
                    if (target.status === 200 && adoptedData.length) {
                        setAdoptedData(adoptedData);
                    }
                })
                .catch(error => {
                    console.log(error);
                    setState({ ...state, errorMessage: error });
                })
        }
    }, [state.cookie]);

    return (
        <>
            {state.backRedirect && <Redirect to={{ pathname: "/animal-control-report", state: { from: `/animal-control-report/${month}-${year}` } }} />}
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Animal Control Drill Down Report</h1>
                <button type="button" onClick={() => setState({ ...state, backRedirect: true })} className="back">Back to Animal Control Report</button>
                <h3>Dogs Surrendered</h3>
                <table>
                    <thead>
                        <td>Dog ID</td>
                        <td>Breed</td>
                        <td>Sex</td>
                        <td>Alteration Status</td>
                        <td>Microchip ID</td>
                        <td>Surrender Date</td>
                    </thead>
                    <tbody>
                        {surrenderedData.length > 0 && (
                            surrenderedData.map(dog => {

                                return (
                                    <tr>
                                        <td>{dog.dogId}</td>
                                        <td>{dog.breed}</td>
                                        <td>{dog.sex}</td>
                                        <td>{dog.alterationStatus ? "Altered" : "Not yet altered"}</td>
                                        <td>{dog.microchipId || "None recorded"}</td>
                                        <td>{dog.surrenderDate}</td>
                                    </tr>
                                )
                            })
                        )}
                    </tbody>
                </table>

                <h3>Dogs Adopted</h3>
                <table>
                    <thead>
                        <td>Dog ID</td>
                        <td>Breed</td>
                        <td>Sex</td>
                        <td>Alteration Status</td>
                        <td>Microchip ID</td>
                        <td>Surrender Date</td>
                        <td>Days in Rescue</td>
                    </thead>
                    <tbody>
                        {adoptedData.length > 0 && (
                            adoptedData.map(dog => {
                                console.log("running")
                                const surrendered = new Date(dog.surrenderDate);
                                const adopted = new Date(dog.adoptionDate);
                                const daysInShelter = Math.round((adopted - surrendered) / (1000 * 60 * 60 * 24));
                                return (
                                    <tr>
                                        <td>{dog.dogId}</td>
                                        <td>{dog.breed}</td>
                                        <td>{dog.sex}</td>
                                        <td>{dog.alterationStatus ? "Altered" : "Not yet altered"}</td>
                                        <td>{dog.microchipId || "None recorded"}</td>
                                        <td>{dog.surrenderDate}</td>
                                        <td>{daysInShelter}</td>
                                    </tr>
                                )
                            })
                        )}
                    </tbody>
                </table>

                {
                    state.errorMessage && (
                        <div className="error-message">{state.errorMessage}</div>
                    )
                }
            </main>
        </>
    );
}

export default AnimalControlDrillDownReport;