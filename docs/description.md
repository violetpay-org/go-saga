# Saga Framework

[한국어](description-kr.md)

- [Saga Framework](#saga-framework)
    * [Overview](#overview)
    * [Components](#components)
        + [Channel](#channel)
        + [MessageRelayer](#messagerelayer)
        + [ChannelRegistry](#channelregistry)
        + [Definition](#definition)
        + [Step](#step)
        + [Endpoint](#endpoint)
        + [Saga](#saga)
        + [Session](#session)
        + [Registry](#registry)
    * [Overall Flow](#overall-flow)

## Overview

The Saga Framework is utilized to ensure transactions between microservices. Both frameworks share the same logic and domain.

- Typescript: [violetpay-org/point3-typescript-saga](https://github.com/violetpay-org/point3-typescript-saga)

- Go: [violetpay-org/go-saga](https://github.com/violetpay-org/go-saga)


## Components

### Channel

There are two types of channels. Each should be used according to the situation.

- \<To Registry\> Side

This is the gateway that sends messages to the Registry so that the Saga can consume them. These messages represent the success/failure responses of other services. The `Send()` method of this channel is responsible for sending messages to the Registry.

- \<MessageRelayer\> Side

This channel serves as a gateway that allows the Saga to publish messages to other services. Messages are not published immediately but are stored in the repository of the channel to be sent by the MessageRelayer. The `Send()` method, which is called by the MessageRelayer, needs to be implemented.

### MessageRelayer

The MessageRelayer reads messages from the repository of channels registered in the ChannelRegistry and publishes them using the channel’s `Send()` method. The `Send()` method must be implemented by the user of the framework. This process can be executed in batches at regular intervals (e.g., every second).

### ChannelRegistry

The ChannelRegistry is a registry that holds multiple channels. Various channels can be registered for use by the MessageRelayer.

### Definition

A Definition outlines the structure of a Saga and can contain multiple steps.

### Step

A Step represents a unit of work that is executed within a single service.

1. Each Step has a unique endpoint.
2. Steps include invocations configured through the endpoint.
3. Steps can execute tasks locally or via external (remote) services.
4. If a subsequent step fails, the preceding steps may trigger compensations (rollbacks).
5. A step can define compensation tasks that will be executed to roll back prior actions in case of a failure. Compensation is defined similarly to invocation, using an endpoint.

### Endpoint

An endpoint defines the actual work to be performed.

- Endpoint (Remote Endpoint)

An Endpoint defines work to be performed by an external service. It includes a message (Command) and a CommandConstructor needed to generate the command when the endpoint is initialized. When the RemoteStep starts:

1. The Orchestrator generates a command using the CommandConstructor.
2. The generated command is stored in the \<MessageRelayer\> Channel repository.

- LocalEndpoint

A LocalEndpoint defines work to be performed locally. The `handle()` method implements the actual task. After the task is executed, a response message can be generated to indicate success or failure (which can be referenced by other external services). Therefore, a channel and response constructors (SuccessResponseConstructor / FailureResponseConstructor) are provided to allow for this:

1. The Orchestrator executes the `handle()` method, then uses the constructor to generate either a success or failure response based on the result.
2. The generated response is stored in the Channel repository.

### Saga

A Saga is a unit of work defined by a Definition. When a saga is executed, a session is created, and each step can share data through the session. A single session corresponds to a single saga execution.

### Session

A session allows for data exchange between multiple steps. It holds the CurrentStep and State. The State can be one of `Compensating`, `Failed`, `Completed`, or `Retrying`.

When a saga is executed, a session is created, which remains active until the saga execution is completed or fails. Once the saga ends, the session is considered a Dead Session, indicating that its state is either `Failed` or `Completed`.

### Registry

The Registry is a collection of multiple sagas. Various sagas can be registered for use within the Registry.

## Overall Flow

![Overall Flow Image](saga.png)

1. The process starts with `Registry.StartSaga()`.
2. It then follows the flow of ***Saga → Definition → Step → Endpoint***, where a message is stored in the Channel repository. Depending on the type of endpoint, the message can be a command or a response.
3. While the message has not yet been published, the saga remains in a pending state, waiting for a response.
4. The MessageRelayer publishes the message:
    1. For local endpoints, the message stored in the repository is a response, and the MessageRelayer sends it through the channel.
    2. For remote endpoints, the message stored in the repository is a command, which the MessageRelayer sends to the external service through the channel.
        1. The external service generates a success/failure response. This response should be delivered to the \<To Registry\> channel through a designated path. In the example above, a message queue is used as a pipeline, and the message reaches the channel through a consumer (which must be designed and implemented by the user).
5. The response that reaches the \<To Registry\> channel is consumed by the Registry. The Orchestrator completes the step based on the response and determines whether to proceed to the next step or trigger compensation.
