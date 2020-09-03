import React, { useEffect, useReducer, useState } from "react";
import { useParams, Redirect } from "react-router-dom";

import { getData, sendData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const DogDetails = () => {

    const [pageState, setPageState] = useState({
        loginRedirect: false,
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer !== null) {
            setPageState({ ...pageState, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    const { dogId } = useParams();

    const reducer = (state, action) => {
        switch (action.type) {
            case "setDog":
                return { ...state, ...action.payload };
            case "sex":
                return { ...state, sex: action.payload };
            case "alterationStatus":
                return { ...state, alterationStatus: action.payload === "true" };
            case "microchipId":
                return { ...state, microchipId: action.payload };
            case "expenses":
                return { ...state, expenses: action.payload };
            case "availableBreeds":
                return { ...state, availableBreeds: action.payload };
            case "errorMessage":
                return { ...state, errorMessage: action.payload };
            case "showModal":
                return { ...state, id: action.payload, showModal: true };
            case "dogDashboard":
                return { ...state, dashboardRedirect: true };
            case "breed":
                const newBreed = (() => {
                    if (action.payload.checked) {
                        if (action.payload.value === "Mixed" || action.payload.value === "Unknown") {
                            const falseBreeds = Object.keys(state.breed).reduce((breeds, breed) => {
                                breeds[breed] = false;
                                return breeds;
                            }, {});
                            return { ...falseBreeds, [action.payload.value]: true };
                        } else {
                            return { ...state.breed, [action.payload.value]: true };
                        }
                    } else {
                        return { ...state.breed, [action.payload.value]: false }
                    }
                })();
                return { ...state, breed: newBreed };
            case "canEditBreeds":
                return { ...state, canEditBreeds: true };
            case "canEditSex":
                return { ...state, canEditSex: true };
            case "canEditAlterationStatus":
                return { ...state, canEditAlterationStatus: true };
            case "canEditMicrochipId":
                return { ...state, canEditMicrochipId: true };
            case "showAddExpenseForm":
                return { ...state, showAddExpenseForm: action.payload };
            case "newExpenseVendor":
                return { ...state, newExpense: { ...(state.newExpense), vendor: action.payload } };
            case "newExpenseDate":
                return { ...state, newExpense: { ...(state.newExpense), date: action.payload } };
            case "newExpenseAmount":
                if (typeof +action.payload !== "number" || +action.payload < 0) {
                    return { ...state, expenseErrorMessage: "Please enter a valid numeric amount." }
                } else {
                    if (state.expenseErrorMessage === "Please enter a valid numeric amount.") {
                        return { ...state, expenseErrorMessage: "" }
                    }
                }
                return { ...state, newExpense: { ...(state.newExpense), amountInCents: Math.floor(+action.payload * 100) } };
            case "newExpenseDescription":
                return { ...state, newExpense: { ...(state.newExpense), description: action.payload } };
            case "resetNewExpense":
                return { ...state, newExpense: { dogIdFk: +dogId } };
            case "addExpense":
                return { ...state, expenses: [...(state.expenses), action.payload] };
            case "expenseErrorMessage":
                return { ...state, expenseErrorMessage: action.payload }
            case "addAdoptionRedirect":
                return { ...state, addAdoptionRedirect: true };
            case "isAdoptable":
                return { ...state, isAdoptable: true };
            default:
                throw new Error();
        }
    };

    const [dog, dispatch] = useReducer(
        reducer,
        {
            name: "",
            breed: {},
            sex: null,
            alterationStatus: null,
            dateOfBirth: "",
            description: "",
            microchipId: "",
            surrenderDate: "",
            surrenderReason: "",
            surrenderedByAnimalControl: null,
            volunteerId: -1,
            expenses: [],
            errorMessage: "",
            showModal: false,
            canEditBreeds: false,
            canEditSex: false,
            canEditAlterationStatus: false,
            canEditMicrochipId: false,
            showAddExpenseForm: false,
            newExpense: { dogIdFk: +dogId },
            expenseErrorMessage: "",
            addAdoptionRedirect: false,
            isAdoptable: false
        });

    useEffect(() => {
        if (pageState.cookie.hasOwnProperty("isTrustedVolunteer") && pageState.cookie.isTrustedVolunteer !== null) {
            const breedsObj = {};
            getData("/breed")
                .then(({ target }) => {
                    JSON.parse(target.response).forEach(breed => {
                        breedsObj[breed.breed] = false;
                    });
                }).then(() => {
                    getData(`/dog/${dogId}`)
                        .then(({ target }) => {
                            const dog = JSON.parse(target.response);
                            dog.breed.forEach(breed => {
                                breedsObj[breed] = true;
                            });
                            dog.breed = breedsObj;

                            // derive permissions
                            if (dog.breed["Unknown"] || dog.breed["Mixed"]) {
                                dispatch({ type: "canEditBreeds" });
                            }
                            if (dog.sex === "unknown") {
                                dispatch({ type: "canEditSex" });
                            }
                            if (!dog.alterationStatus) {
                                dispatch({ type: "canEditAlterationStatus" });
                            }
                            dog.microchipId = dog.microchipId.String;
                            if (!dog.microchipId) {
                                dispatch({ type: "canEditMicrochipId" });
                            }

                            if (dog.microchipId && dog.alterationStatus) {
                                dispatch({ type: "isAdoptable" });
                            }

                            // set dog
                            dispatch({ type: "setDog", payload: dog });
                        })
                })
                .catch(error => {
                    dispatch({ type: "errorMessage", payload: error });
                    console.log("error calling the API:", error);
                })
        }
    }, [dogId, pageState.cookie]);

    const updateDog = () => {
        const dogData = {
            breed: Object.keys(dog.breed).filter(breed => {
                return dog.breed[breed];
            }),
            sex: dog.sex,
            alterationStatus: dog.alterationStatus
        };
        dogData.microchipId = { "String": dog.microchipId, "Valid": dog.microchipId.length ? true : false };
        if (!Object.keys(dog.breed).length) {
            dispatch({ type: "errorMessage", payload: "Please select at least one dog breed." });
            return;
        }
        sendData(dogData, `/dog/${dogId}`, "put")
            .then(({ target }) => {
                if (target.status === 204) {
                    // show success modal
                    /* eslint-disable no-restricted-globals */
                    location.reload();
                } else {
                    console.log("error adding dog (non-200):", target.response);
                    dispatch({ type: "errorMessage", payload: target.response });
                }
            })
            .catch(error => {
                // show error modal
                console.log("error adding dog (catch):", error);
                dispatch({ type: "errorMessage", payload: JSON.parse(error) });
            });
    }

    const addExpense = () => {
        sendData(dog.newExpense, `/expense`)
            .then(({ target }) => {
                if (target.status === 201) {
                    dispatch({ type: "resetNewExpense" });
                    dispatch({ type: "showAddExpenseForm", payload: false });
                    dispatch({ type: "addExpense", payload: JSON.parse(target.response) });
                    dispatch({ type: "expenseErrorMessage", payload: "" });
                } else {
                    console.log(target.response);
                    dispatch({ type: "expenseErrorMessage", payload: target.response });
                }
            })
            .catch(error => {
                console.log(error);
                dispatch({ type: "expenseErrorMessage", payload: error });
            });
    }

    const Checkbox = ({ breed }) => {
        return (
            <p key={breed}>
                <input
                    type="checkbox"
                    value={breed}
                    name={breed.split(" ").join("-")}
                    checked={dog.breed[breed]}
                    onChange={({ target }) => {
                        dispatch({ type: 'breed', payload: target })
                    }}
                    disabled={(dog.breed["Mixed"] || dog.breed["Unknown"]) && !(breed === "Mixed" || breed === "Unknown")}
                /><label htmlFor={breed} style={{ display: "inline", paddingLeft: "10px" }}>{breed}</label>
            </p>
        )
    }

    return (
        <>
            {dog.dashboardRedirect && <Redirect to={{ pathname: `/dog-dashboard` }} />}
            {dog.addAdoptionRedirect && <Redirect to={{ pathname: `/add-adoption`, state: { dog } }} />}
            <NavLinks isTrustedVolunteer={pageState.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Dog Details</h1>
                <button type="button" onClick={() => dispatch({ type: "dogDashboard" })} className="back">Back to Dog Dashboard</button>
                {dog.isAdoptable && pageState.cookie.isTrustedVolunteer &&
                    <button type="button" onClick={() => dispatch({ type: "addAdoptionRedirect" })} className="back">Add Adoption</button>}
                <form>
                    <div className="form-unit">
                        <span>Name: {dog.name}</span>
                    </div>
                    {dog.canEditBreeds ? (
                        <div className="form-unit">
                            <label htmlFor="breed">Dog Breed: (Please select at least one.)</label>
                            {Object.keys(dog.breed).map(b => {
                                return <Checkbox breed={b} key={b} />
                            })}
                        </div>
                    ) : (
                            <div className="form-unit">
                                <span>Breed: {Object.keys(dog.breed).filter(breed => {
                                    return dog.breed[breed];
                                }).join("/")}</span>
                            </div>
                        )}
                    <div className="form-unit">
                        {dog.canEditSex ? (
                            <>
                                <label htmlFor="sex">Dog Sex:</label>
                                <p>
                                    <input type="radio"
                                        id="male"
                                        name="sex"
                                        value="male"
                                        onChange={({ target }) => {
                                            dispatch({ type: "sex", payload: target.value })
                                        }}
                                        checked={dog.sex === "male"}
                                    />
                                    <label htmlFor="male" style={{ display: "inline" }}>Male</label>
                                </p>
                                <p>
                                    <input type="radio"
                                        id="female"
                                        name="sex"
                                        value="female"
                                        onChange={({ target }) => {
                                            dispatch({ type: "sex", payload: target.value })
                                        }}
                                        checked={dog.sex === "female"}
                                    />
                                    <label htmlFor="female" style={{ display: "inline" }}>Female</label>
                                </p>
                                <p>
                                    <input type="radio"
                                        id="unknown"
                                        name="sex"
                                        value="unknown"
                                        onChange={({ target }) => {
                                            dispatch({ type: 'sex', payload: target.value })
                                        }}
                                        checked={dog.sex === "unknown"}
                                    />
                                    <label htmlFor="unknown" style={{ display: "inline" }}>Unknown</label>
                                </p>
                            </>
                        ) : (
                                <span>Sex: {dog.sex}</span>
                            )}
                    </div>
                    <div className="form-unit">
                        {dog.canEditAlterationStatus ? (
                            <>
                                <label htmlFor="alterationStatus">Dog Alteration Status:</label>
                                <p>
                                    <input type="radio"
                                        id="altered"
                                        name="alterationStatus"
                                        value={true}
                                        onChange={({ target }) => {
                                            dispatch({ type: 'alterationStatus', payload: target.value })
                                        }}
                                        checked={dog.alterationStatus}
                                    />
                                    <label htmlFor="altered" style={{ display: "inline" }}>Spayed/Neutered</label>
                                </p>
                                <p>
                                    <input type="radio"
                                        id="unaltered"
                                        name="alterationStatus"
                                        value={false}
                                        onChange={({ target }) => {
                                            dispatch({ type: 'alterationStatus', payload: target.value })
                                        }}
                                        checked={!dog.alterationStatus}
                                    />
                                    <label htmlFor="unaltered" style={{ display: "inline" }}>Unaltered</label>
                                </p>
                            </>
                        ) : (
                                <span>Alteration Status: {dog.alterationStatus ? "Spayed/Neutered" : "Unaltered"}</span>
                            )}
                    </div>
                    <div className="form-unit">
                        <span>Date of Birth: {new Date(dog.dateOfBirth).toDateString()}</span>
                    </div>
                    <div className="form-unit">
                        <span>Description: {dog.description}</span>
                    </div>
                    <div className="form-unit">
                        {dog.canEditMicrochipId ? (
                            <>
                                <label>Microchip ID:</label>
                                <input
                                    type="text"
                                    onChange={({ target }) => {
                                        dispatch({ type: "microchipId", payload: target.value })
                                    }} />
                            </>
                        ) : (
                                <span>Microchip ID: {dog.microchipId}</span>
                            )}
                    </div>
                    <div className="form-unit">
                        <span>Surrender Date: {new Date(dog.surrenderDate).toDateString()}</span>
                    </div>
                    <div className="form-unit">
                        <span>Surrender Reason: {dog.surrenderReason}</span>
                    </div>
                    <div className="form-unit">
                        <span>Surrender was by Animal Control: {dog.surrenderedByAnimalControl ? "Yes" : "No"}</span>
                    </div>
                    <div className="form-unit">
                        <span>Surrender recorded by volunteer {dog.volunteerId}</span>
                    </div>
                    {(
                        dog.canEditBreeds || dog.canEditSex ||
                        dog.canEditAlterationStatus || dog.canEditMicrochipId
                    ) && (
                            <>
                                <button type="button" onClick={updateDog}>Update Dog</button>
                                <span style={{ color: "red", paddingLeft: "10px" }}>{dog.errorMessage}</span>
                            </>
                        )}
                </form>
                <button type="button" onClick={() => dispatch({ type: "dogDashboard" })}> Back to Dog Dashboard</button >
                <hr></hr>
                <h1>Expenses</h1>
                <h3>Dog Total : ${(dog.expenses.reduce((totalExpense, currentExpense) => totalExpense += currentExpense.amountInCents, 0)) / 100}</h3>
                {dog.showAddExpenseForm ? (
                    <div>
                        <h3>New Expense</h3>
                        <form>
                            <div className="form-unit">
                                <label id="expense-date">Date</label>
                                <input
                                    type="date"
                                    onChange={
                                        ({ target }) => dispatch({ type: "newExpenseDate", payload: target.value })
                                    } />
                            </div>
                            <div className="form-unit">
                                <label id="expense-vendor">Vendor</label>
                                <input
                                    type="text"
                                    onChange={
                                        ({ target }) => dispatch({ type: "newExpenseVendor", payload: target.value })
                                    } />
                            </div>
                            <div className="form-unit">
                                <label id="expense-amount">Amount (in dollars)</label>
                            $<input
                                    type="number"
                                    step={0.01}
                                    min={0}
                                    onChange={({ target }) => {
                                        if (typeof +target.value === "number") {
                                            dispatch({ type: "newExpenseAmount", payload: target.value });
                                        }
                                    }
                                    } />
                            </div>
                            <div className="form-unit">
                                <label id="expense-description">Description</label>
                                <input
                                    type="text"
                                    onChange={
                                        ({ target }) => dispatch({ type: "newExpenseDescription", payload: target.value })
                                    } />
                            </div>
                            <button type="button" onClick={addExpense}>Submit Expense</button>
                            {dog.expenseErrorMessage && <span style={{ color: "red", paddingLeft: "10px" }}>{dog.expenseErrorMessage}</span>}
                        </form>
                    </div>
                ) : (
                        <>
                            <button type="button" onClick={() => dispatch({ type: "showAddExpenseForm", payload: true })}>Add New Expense</button>
                            {dog.expenseErrorMessage && <span style={{ color: "red", paddingLeft: "10px" }}>{dog.expenseErrorMessage}</span>}
                        </>
                    )}

                {
                    dog.expenses.length ? (
                        dog.expenses.map((expense, index) => (
                            <div key={index} className="form-unit">
                                <div>
                                    <span>Date: {new Date(expense.date).toDateString()}</span>
                                </div>
                                <div>
                                    <span>Vendor: {expense.vendor}</span>
                                </div>
                                <div>
                                    <span>Description: {expense.description}</span>
                                </div>
                                <div>
                                    <span>Cost: ${expense.amountInCents / 100}</span>
                                </div>
                            </div>
                        ))
                    ) : (<div>No expenses entered yet.</div>)
                }
            </main>
        </>
    );
}

export default DogDetails;