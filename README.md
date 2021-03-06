# EQ05 - GamesAPI
Ceci est un projet d'exploration afin d'expérimenter la langage Go et le développement "Back-end".
Le thème principal est le jeu vidéo. La vision est de créer une interface REST évolutive qui offrira plusieurs fonctionnalités sur différents axes du développement logiciel afin de pouvoir faire la gestion de votre collection de jeux, contenus téléchargeables, listes d'amis, et plus encore!
## Sécurité
Cette interface pourra ultimement intégrer des fonctionnalités communes en sécurité, comme la communication par TLS, l'authentification, l'autorisation par rôle, et plus encore au fil du projet.
## Évolutivité
Ce projet contiendra notamment de la documentation Docker/Kubernetes, afin de facilement générer et déployer la solution sur le nuage.
## Pertinence
Cette interface pourra ultimement faire le pont avec d'autre plateformes populaires comme Steam ou Battle.Net (Blizzard).
## Accessibilité
Cette interface devra respecter certains standards afin d'être familière aux développeurs qui se l'approprieront.
# Mise en route
Les prérequis pour compiler et exécuter le projet sont les suivants: 
- Docker Engine (se nomme 'Docker for Desktop' sous Windows)
- image Docker pour MSSQL

## Base de données
Si vous n'avez pas d'image Docker MSSQL, voici une marche à suivre:
1. Ouvrir un invite de commande dans ./dockerizedDatabase
2. Exécuter `docker build --tag sqlserver_image:latest .`
3. Vous avez désormais une image MSSQL valide pour l'environnement.


## Environnement de développement
Marche à suivre pour lancer un serveur Dev avec base de données MSSQL et mise à jour automatique:
1. Ouvrir un invite de commande à la racine du projet
2. Taper `docker-compose up development db`
3. Codez comme s'il n'y avait pas de lendemain!

## Documentation

Marche à suivre pour générer la documentation: 
1. [Télécharger et installer `nodejs`](https://nodejs.org/en/) (Installera `npm`)
2. Installer `raml2html` via la commande `npm install raml2html -g`
3. Dans une nouvelle console, exécuter `raml2html Documentation/Raml/main.raml > doc.html`


