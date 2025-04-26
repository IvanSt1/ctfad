import axios from 'axios';

export async function sendCard(_to, data, cardType) {
    try {
        const resp = await axios.post('http://localhost:8083/api/add_card', {
            to: _to,
            text: data,
            imageType: cardType
        }, {
                withCredentials: true,
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
        });
        if (resp.status !== 200) {
            throw new Error(resp.data.toString());
        }
    } catch(error) {
        throw new Error("Could not add card");
    }
}

export async function getCard(cardId) {
    try {
        const resp = await axios.get(`http://localhost:8083/api/card`,{
            params: {
                id: cardId
            },
                withCredentials: true,
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            });
        return resp.data;
    } catch(reason) {
        throw new Error(reason.message)
    }
}

export async function getLatestCards(number = "") {
    try {
        const resp = await axios.get(`http://localhost:8083/api/cards`,{
            params: { latest: number.toString() },
                withCredentials: true,
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            });
        return resp.data;
    } catch(reason) {
        throw new Error(reason.message);

    }
}
