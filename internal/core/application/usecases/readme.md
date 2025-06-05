# Definition

Use cases Have names that reflect user intentions eg. GenerateAssessment, ExportAssessment, SubcscribeToPlan. Services can and are often used within use cases, Use cases ochestrates aggregrates and services. So you could say use case is higher up as it can contain services

They are on per user intention, while services will probably support multiple usecases

---

🧠 Think of it this way:
Use Case: "What does the user want to do?"

Service: "What logic or behavior supports that goal?"

Domain services exist within the domain and should be side effect free, use multiple entities and requires domain rules

---

🧠 How Domain Services Relate to Aggregates
Aggregates manage consistency and state across entities

Domain services manage domain logic across aggregates/entities
