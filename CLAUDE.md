# CLAUD.md

Collaborative Live Application for Unified Documents

## Overview

This project aims to build a Proof of Concept (PoC) for a real-time document editing application.

## Features

- Document creation
- User authentication
- Real-time collaborative document editing

## Screen Flow

- There are three main screens:
  - Login screen
  - Document list page
  - Collaborative editing screen
- After logging in, the document list is displayed.
- Users can open a document from the list and start editing it.

## Technologies

- **Frontend**
  - TypeScript
  - Vue

- **Backend**
  - Go
  - net/http

- **Communication**
  - WebSocket
  - If possible, modern protocols such as QUIC or WebTransport should be used.

## Notes and Future Plans

- Ultimately, the app should be usable as a PukiWiki plugin.
- The server will be implemented separately, while the client will be implemented as a PukiWiki plugin.
- The UI/UX for PukiWiki integration is undecided.
  - It does not need to be implemented at this stage.
  - However, this context should be kept in mind as advice on integration may be sought in the future.
