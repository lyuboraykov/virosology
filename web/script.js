(function(){
    document.getElementById("start-button").addEventListener("click", main);
    const canvas = document.getElementById("view");
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight - 100;

    function main() {
        const populationSize = parseInt(document.getElementById("pop-size").value),
              initialInfectedCount = parseInt(document.getElementById("init-infected").value),
              chanceOfTransmission = parseFloat(document.getElementById("trans-chance").value) / 100,
              chanceOfDeath = parseFloat(document.getElementById("death-chance").value) / 100,
              daysUntilDeath = parseInt(document.getElementById("days-to-death").value),
              isolationLevel = parseFloat(document.getElementById("isolation-percentage").value) / 100,
              daysToRecover = parseInt(document.getElementById("recover-days").value),
              intervalBetweenFrames = 100; // ms

        document.getElementById("content").classList.add("hidden");

        const canvas = document.getElementById('view');
        canvas.classList.remove("hidden");
        const ctx = canvas.getContext('2d');
        ctx.font = '30px serif';

        const population = newPopulation(populationSize, initialInfectedCount, daysToRecover, isolationLevel, canvas.width, canvas.height);
        let currentDay = 1;
        let interval = setInterval(() => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            let infectedCount = 0;
            for (const p of population) {
                if (p.isInfected(currentDay)) {
                    infectedCount++;
                }
            }
            if (infectedCount == 0) {
                clearInterval(interval);
                hideCanvas();
                displayStats(population, canvas);
            }

            for (const p of population) {
                let personEmoji = '🙂';
                if (p.isInfected(currentDay)) {
                    personEmoji = '🥵';
                } else if (p.isImmune(currentDay)) {
                    personEmoji = '🥳';
                } else if (!p.isAlive) {
                    personEmoji = '☠️';
                }

                ctx.fillText(personEmoji, p.position.x, p.position.y);
            }
            movePopulation(population, chanceOfTransmission, chanceOfDeath, canvas.width, canvas.height, currentDay, daysUntilDeath);
            currentDay++;
        }, intervalBetweenFrames);
        const stopButton = document.getElementById("stop");
        stopButton.classList.remove("hidden");
        stopButton.addEventListener("click", () => {
            clearInterval(interval);
            hideCanvas();
        });
    }

    function displayStats(population, canvas) {
        let totalInfected = 0, totalDead = 0;
        for (p of population) {
            if (p.infectedAt != 0) {
                totalInfected++;
            }
            if (!p.isAlive) {
                totalDead++;
            }
        }
        document.getElementById("stats").innerText = `Total infected: ${totalInfected}; total dead: ${totalDead}`;
    }

    function hideCanvas() {
        canvas.classList.add("hidden");
        document.getElementById("content").classList.remove("hidden");
        document.getElementById("stop").classList.add("hidden");
    }

    class Position {
        constructor(x, y) {
            this.x = x;
            this.y = y;
        }
    }

    function randomToN(n) {
        return (Math.random() * n) | 0;
    }

    class Person {
        constructor(position, infectedAt, isIsolated, isAlive, daysToRecover) {
            this.position = position;
            this.infectedAt = infectedAt;
            this.isIsolated = isIsolated;
            this.isAlive = isAlive;
            this.daysToRecover = daysToRecover;
        }

        isInfected(currentDay) {
            return this.isAlive && this.infectedAt != 0 &&
                currentDay - this.infectedAt < this.daysToRecover;
        }

        isImmune(currentDay) {
            return this.isAlive && this.infectedAt != 0 && currentDay - this.infectedAt >= this.daysToRecover;
        }

        daysInfected(currentDay) {
            return currentDay - this.infectedAt;
        }
    }

    function newPopulation(populationSize,
        initialInfectedCount,
        daysToRecover,
        isolationLevel,
        maxX,
        maxY
    ) {
        population = [];
        for (i = 0; i < populationSize; i++) {
            while (true) {
                let x = randomToN(maxX);
                let y = randomToN(maxY);
                let candidatePosition = new Position(x, y);
                let takenBy = positionTakenBy(candidatePosition, population);
                if (takenBy == null) {
                    let infectedAt = 0;
                    if (i < initialInfectedCount) {
                        infectedAt = 1;
                    }
                    population.push(new Person(
                        candidatePosition,
                        infectedAt,
                        wouldOccurWithChance(isolationLevel),
                        true,
                        daysToRecover,
                    ))
                    break;
                }
            }

        }

        return population;
    }

    function positionTakenBy(position, population) {
        for (const p of population) {

            const xInBorders = (p.position.x - 10 < position.x && position.x < p.position.x + 10);
            const yInBorders = (p.position.y - 10 < position.y && position.y < p.position.y + 10);

            if (xInBorders && yInBorders) {
                return p;
            }
        }
        return null;
    }

    function wouldOccurWithChance(chance) {
        return randomToN(100) < chance * 100;
    }

    function movePopulation(population, chanceOfTransmission, chanceOfDeath, maxX, maxY, currentDay, daysUntilDeath) {
        for (const p of population) {
            if (p.isIsolated || !p.isAlive) {
                continue;
            }
            if (p.isInfected(currentDay) && wouldOccurWithChance(chanceOfDeath) && p.daysInfected(currentDay) == daysUntilDeath) {
                p.isAlive = false;
                continue;
            }
            const xOrY = wouldOccurWithChance(0.5);
            const minusOrPlus = wouldOccurWithChance(0.5);
            let direction = 1;
            if (minusOrPlus) {
                direction = -1;
            }
            const candidatePosition = new Position(p.position.x, p.position.y);
            if (xOrY) {
                candidatePosition.x += direction * 20;
            } else {
                candidatePosition.y += direction * 20;
            }

            if (candidatePosition.x < 0 || candidatePosition.y < 0 || candidatePosition.x >= maxX || candidatePosition.y >= maxY) {
                continue;
            }

            let takenBy = positionTakenBy(candidatePosition, population);
            if (takenBy != null) {
                if (!takenBy.isAlive) {
                    continue;
                }
                if (p.isInfected(currentDay) && !takenBy.isInfected(currentDay) && !takenBy.isImmune(currentDay)) {
                    if (wouldOccurWithChance(chanceOfTransmission)) {
                        takenBy.infectedAt = currentDay;
                    }
                    continue;
                }
                if (takenBy.isInfected(currentDay) && !p.isInfected(currentDay) && !p.isImmune(currentDay)) {
                    p.infectedAt = currentDay;
                }
                continue;
            }

            p.position = candidatePosition;
        }
    }
})()
