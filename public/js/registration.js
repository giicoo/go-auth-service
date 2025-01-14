const url = "http://localhost:8080/create-user";

// Данные для отправки
const data = {
    email: "dad",
    password: "123"
};

// Функция для выполнения POST-запроса
async function sendPostRequest(url, data) {
    try {
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data) // Преобразуем данные в JSON-строку
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        console.log("Server response:", result);
    } catch (error) {
        console.error("Error:", error);
    }
}

// Выполнить запрос
sendPostRequest(url, data);
