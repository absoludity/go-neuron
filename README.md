Go Neuron
=========

go-neuron is a personal experimental library for simulating neurons and
neural networks.

A neuron implements a simple interface that responds to additional
potential being added at a point in time. If the additional potential
over time reaches a threshold, the neuron fires and sends an
activation event on a channel for processing.

[![Build Status](https://travis-ci.org/absoludity/go-neuron.png)](https://travis-ci.org/absoludity/go-neuron)


Action Potentials
-----------------

The action-potential interface enables getting the potential (at a specific
point in time) as well as adding to the potential (at a specific point in time).
Currently there is only one action potential implementation, the SimpleActionPotential.


Neurons
-------

A Neuron is a composition of an action potential and an Axon which carries the
signal to the terminals (connecting other neurons) with a specified propagation delay.
