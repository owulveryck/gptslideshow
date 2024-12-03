## Context

I have been an early advocate for the data mesh paradigm since Zhamak Dehghani first proposed it.
As an early supporter, I identified the potential of this novel approach to data organization.
Four years have passed, and the data mesh paradigm has indeed earned widespread acceptance.
However, I have not yet seen a concise, practical data mesh transformation plan within any organization.

When I say "advocate," I mean to say that I've been highlighting the benefits of this paradigm, which are rooted in its four main pillars:

- The orientation of design towards **domains**
- Applying **product thinking to data** (data-as-a-product)
- Federated computational **governance**
- Developing self-service **data platforms**

Then comes the question, **where should one actually start to implement the data mesh?**

In principle, any organization can kick off their journey towards data mesh by giving focus to these four pillars.

**Beginning with** a design that is **domain-oriented** builds the groundwork for a thorough comprehension of the data mesh.
This means not only setting the data mesh as an objective but also **ensuring that the decomposition** of the **domain syncs with** the **existing structure** of the organization.
However, this is a **profoundly conceptual approach** that might not yield immediate results, and moreover, it lacks the agility that is so beneficial.

Both **federated computational governance** and the **self-service data platform** are simply **enablers** of the data mesh.
They share a common objective: to simplify the development of data-as-a-product and the **creation of interconnections**, essentially supporting the mesh.
One can try implementing them as a foundation, but to mesh what ?

What remains then is to tackle data-as-a-product, a cornerstone of the data mesh that I have previously discussed.

Interestingly, several organizations claim to have implemented the **data mesh "by accident,"** perceiving this paradigm as the natural evolution of data management.

In this article, I attempt to apply a well-recognized model of evolutionary progression to understand data evolution.
**The objective** is to aid in **visualizing data maturity** and assist companies in **identifying** their **tipping point**,
i.e., when they will start seeing **significant benefits** from **implementing data contracts** and **treating data as a product**.

## Modeling the evolution

I will first explain the model I will use.
This model is known as the evolution model by Simon Wardley and is successfully implemented in Wardley Maps.
My goal here is not to describe a specific company's landscape, so I won't need a full map.
Instead, I will use the evolution model and try to apply its general purpose function to the data.

**Disclaimer: Regarding the model:**
The theory of evolution is well-suited for application in a competitive environment where everything evolves based on supply and demand. I am considering businesses that are subject to those constraints of competition and, consequently, their data will also follow those rules. Therefore, the model will apply.

**Why model the evolution?** Understanding evolution is about understanding how components change over time. Modeling evolution is about finding a pattern to potentially provide insights into the future trajectories of those components.

### The model in a glimpse

Simon Warldey needed a way to represent the evolution of the components on his map.
He could not rely on a basic timescale because it would prevent him from comparing heterogeneous elements and would break the consistency of movement.

For example, on a timescale, the distance between the genesis and the maturity of a car (around 100 years) would have been much greater than the distance between the same points for a smartphone (around 10 years).
He discovered that the **evolution** is a **function** of its **ubiquity** **and** its **certainty**.

In a market economy, **ubiquity is led by demand**. More demand induces more presence. It is a declination of the theory of the [diffusion of innovations of Rogger](https://en.wikipedia.org/wiki/Diffusion_of_innovations).
The certainty comes from the Stacey matrix. The matrix postulates that as the availability of key compoents or information increases, the certainty regarding decision-making outcomes also increases, allowing for more predicatble and informed planning and execution.
In a sense, **certainty is driven by supply**.

For example, let's consider a manufacturing company that produces electronic gadgets. In this scenario, one of the critical components they need are semiconductor chips. When the supply of these chips is low due to market shortages or logistic issues, the company faces high uncertainty regarding its production schedules and ability to meet customer demand.

An empirical analysis led to this representation:
![](https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEgatdAD8t3Jp7BEjlcpxMwwUMGPbmu-zs9kwEX4KlVqZ31VwHzShmyAr1ZE0zC4YWUnTXWncgIVFPr6_-CQhKn8FO2He4qs-KGd5CrlLcW7S-ZzNxUZLAQqDQE-Vqe11g8Rt7eOaA/s1600/Screen+Shot+2014-03-15+at+18.48.03.png)
(source [Simon Wardley's blog.garvediance.org](https://blog.gardeviance.org/2014/03/on-mapping-and-evolution-axis.html))

The model is a kind of S curve.

### The model of the data according to Wardley

The analysis of the model allowed for formalizing four steps of evolution labeled by default _Genesis, Custom built, Product, and Commodity_:

![](https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEjMFN3o1ujMDfd4y78hHCRFmPSTf9BP5C_Ej1jtEyZrmNC21aBw-18gAbVk88nKHdVa3gd_-D3z3pKKfO4Wa6XsIa1BuTkeiazqGLdu8vlUPsSaXeDgbkbvrMy3CSHlUiqk5ol1ig/s1600/Screen+Shot+2014-01-09+at+13.26.48.png)
(source [Simon Wardley's blog.garvediance.org](https://blog.gardeviance.org/2014/03/on-mapping-and-evolution-axis.html))

Those are just common labels for a form of capital. For the data, according to Wardley's theory, the labels of the four stages are: **_Unmodeled, divergent, convergent and modeled_**:
![](https://i0.wp.com/learnwardleymapping.com/wp-content/uploads/2020/01/20200122_124314.jpg?resize=1080%2C466&ssl=1)
(source: [learnwardleymapping.com](https://learnwardleymapping.com/2020/01/22/visualizing-the-interaction-of-evolution-and-data-measurement/))

### Deriving the model

#### Certainty of the data

Let's revisit the certainty mechanism to determine if we can adjust the model to accommodate the evolution of data within a business. I consider the certainty equivalent to the level of confidence in the decision taken based on this data. Here are the labels I will use:

- **Raw data**: In my experience, data begins as raw during the exploratory phase. It lacks ubiquity, residing solely in the database and accessible only through a service and/or an API, essentially a **data product** (a product driven by data).
- **Curated data**: This marks the second stage of data certainty. Data experts come into play to ensure accuracy and relevance of data representation to the business.
- **Authoritative**: The final stage of certainty. Data is relevant, complete, consistent, documented, and endorsed by domain experts.

The **raw** data correspond to the **first stage** of the evolution. This is a stage where we define Proof of concepts for example. Then the **curated** data is linked to phases two and three. And eventually, the **last stage** is when the data is **authoritative**.

#### The labels of the four steps of evolution

Regarding the notion of certainty and ubiquity, let's categorize the 4 stages of evolution:

1. **POC**: This stage involves validating concepts.
2. **Application**: In this stage, the data is neatly linked with a specific use-case.
3. **Domain**: This stage is where it gets interesting: the data represents a solution that can be used to address various use-cases within the same domain (think of the domain as a problem space, similar to in DDD).
4. **Enterprise**: This stage encompasses all domains, representing the total of all problems addressed by a business.

Here is the representation of those elements on a diagram:

### The representation

![An S curve representing the evolution of the data, the X is the certainty and Y is the ubiquity. There is a division: the bottom of the S is raw data, the middle is curated, and the top is authoritative. The inflection point is indicated as data-contract.](/assets/data_certitude.svg)

## Using the diagram: Data-as-a-product and data-contract

Now, let's use the diagram.

The data will likely follow the evolution S curve. What is interesting is that the evolution of the properties of the data.
Turning raw data into curated data is mastered. There are major design processes that are helpful in such a transition.

Turning the curated data into authoritative data implies that the data is accessible and usable, maintained, accurate, but the switch is mainly that the data is **endorsed by trusted parties**.
At the scale of the business, it means that the domain is responsible and accountable for its data as the domain is, by default, a trusted party in the organization regarding a specific business area.

The transition is not so sharp when the data leave its prison: when it is exposed to the domain.

This is the point where the product thinking applied to the data brings value. And this is the point where a data-contract is helpful to:

- Facilitate the integration in other use cases of the domain
- Bring trust in the data

Therefore, thinking of data as a product, like any other product, is something that is required in the exploration stage (it can even be seen as over-engineering), but the model illustrates how it is important to treat the data as a product to serve a general purpose for the business.

## Conclusion

In recapping, I have always grappled with one question: _where does one begin when seeking to implement the data mesh paradigm?_
Through the journey of exploring this concept, my most recent and profound insight is: the most strategic starting point lies with the data product.

The presented model emphasizes the pivotal role of the data product. It is projected as an effective solution to an imperative issue: its significant importance becomes evident as data migrates from a single application sphere into the broader domain.
Beyond this, it becomes absolutely critical when the data is expected to deliver tangible value that surpasses its original, defined domain.

The next phase of our journey in understanding the data mesh paradigm involves formalizing a method to assess data maturity accurately.
By examining each piece of data, contract by contract, and domain by domain, we move closer to building a comprehensive and effective mesh.
Throughout this process, remembering to consider data as a product is crucial. Doing so will reap rewards for an organization as it evolves and matures in its data management strategies.
