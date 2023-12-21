WITH SeatCHK AS (
    SELECT
        id,
        student,
        ROW_NUMBER() OVER (ORDER BY id) AS row_num
    FROM
        Seat
)

SELECT
        CASE
            WHEN row_num % 2 = 0 AND row_num < (SELECT MAX(row_num) FROM SeatCHK) THEN id + 1
        ELSE id
END AS id,
    student
FROM
    SeatCHK
ORDER BY
    id;
