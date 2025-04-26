import axios from 'axios';

import { setGender } from '@/service/genderLib';


async function login(_username, _password) {
    await axios.post('http://localhost:8083/api/login', {
        username: _username,
        password: _password,
    }, { 
        withCredentials: true,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
    })
    .then((response)=> {
        document.cookie = response.headers['Set-Cookie'];
        setGender(response.data['gender'])
    })
    .catch(() => {
        throw new Error('Invalid username or password');
    });
}

async function register(_username, _password, _gender) {
    await axios.post('http://localhost:8083/api/register', {
       username: _username,
       password: _password,
       gender: _gender
    }, {
        withCredentials: true,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
    })
    .then((response)=> {
        if (response.status !== 200) {
            throw new Error("User already exists or invalid params");
        }
        document.cookie = response.headers['Set-Cookie'];
        setGender(_gender.toLowerCase());
    })
    .catch(() => {
        throw new Error("User already exists or invalid params");
    });
}

async function logout() {
    try {
        await axios.post('http://localhost:8083/api/logout', {}, {
            withCredentials: true,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
        })
    }
    finally {
         localStorage.removeItem("authorized");
         document.cookie = "session="
         localStorage.removeItem("gender");
    }
}

export {login, register, logout}
