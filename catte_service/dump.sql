PGDMP         7                w            Catte    11.2    11.2 
               0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                       false                       0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                       false                       0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                       false                       1262    16396    Catte    DATABASE     �   CREATE DATABASE "Catte" WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'English_United States.1252' LC_CTYPE = 'English_United States.1252';
    DROP DATABASE "Catte";
             postgres    false            �            1259    16405    rooms    TABLE     �   CREATE TABLE public.rooms (
    roomid character varying NOT NULL,
    amount bigint,
    numplayer integer,
    host character varying,
    maxplayer bigint
);
    DROP TABLE public.rooms;
       public         postgres    false            �            1259    16397    users    TABLE     6  CREATE TABLE public.users (
    userid character varying NOT NULL,
    username character varying NOT NULL,
    source character varying(50) NOT NULL,
    password character varying,
    lastcheckin date,
    amount bigint,
    user3rdid character varying,
    dateofbirth date,
    image character varying
);
    DROP TABLE public.users;
       public         postgres    false            �
          0    16405    rooms 
   TABLE DATA                     public       postgres    false    197   i	       �
          0    16397    users 
   TABLE DATA                     public       postgres    false    196    
       �
           2606    16412    rooms roomid 
   CONSTRAINT     N   ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT roomid PRIMARY KEY (roomid);
 6   ALTER TABLE ONLY public.rooms DROP CONSTRAINT roomid;
       public         postgres    false    197            �
           2606    16421    users userid 
   CONSTRAINT     N   ALTER TABLE ONLY public.users
    ADD CONSTRAINT userid PRIMARY KEY (userid);
 6   ALTER TABLE ONLY public.users DROP CONSTRAINT userid;
       public         postgres    false    196            �
   �   x���v
Q���W((M��L�+���-V� Q�):
����y%:
y��9���E:
��@���
���B��O�k��������������9��\�ԵÈv���S:�a4�ЀƖ�=b�Ì�vc��� ���M      �
   �  x��[o�0���)�P�N�`'��lBSDK�0z����ɗc��,	+��_B����w��9��-���ǣ���z�'��4Q�m	E�7!�m����@�)��Pu�EY>آ�KEY��:�ڎ��mV=�B7Z-*�F&E��d#V�Ź�&7�3��0 K�1����@"�}��
��j;-�dq��>B�̢<o¡+	>���v���;%O�Y*5��a�����?�����,[�oegI�ѻ��{��z�p=��:��L&M�������pi�� p#ʌB!s%�#!W���B��v�z]�Ǌ��˯�n����b��Z#�DC��R�:�n^6�yaM�B�����Pϥ��À� S78�!Y�U�ᣇDWq����ƘG9��Q,ʸ���:$e��o��V��%��՚��(�\�0>���Ⱦ����w���x�L�3��x�s˝�-�>�DRD�dH��"��@��V{������G�8~u��p|��IK�#�it^\-���k:�L~���:�ٰ?����?��XxB	�����^[��.x�L�?�wf<�v��+rN���a�0����u�O?�1�ig���y��YY�7#�X^���]�U��e�^�Q7���Gq�CsV�\���u����i糱|�!G��Ƕ���|Հ>8��h��     